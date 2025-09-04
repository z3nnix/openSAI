class ChatLanguageModel
  def initialize
    @word_freq = Hash.new(0)
    @word_pairs = Hash.new { |h, k| h[k] = Hash.new(0) }
    @vocabulary = Set.new
    @phrases = []
    @used_words = Set.new
    @last_question_words = Set.new
    
    # Разделяем стоп-слова на категории
    @stop_words = Set.new(%w[
      это тот такой как чтобы хотя если потому поэтому
      очень слишком вполне совершенно абсолютно совершенно
      уже еще только именно даже же ли бы ведь ведь
      сам сама само сами
      весь вся всё все любой любая любое любые каждый каждая каждое каждые
    ])
    
    # Вопросительные слова - не фильтруем их!
    @question_words = Set.new(%w[
      что кто какой какая какое какие чей чья чье чьи
      где куда откуда когда почему зачем сколько насколько
      как
    ])
    
    # Слова-связки, которые можно фильтровать в некоторых случаях
    @filler_words = Set.new(%w[
      ну вот это же ли бы и в на за к с у о
    ])
  end

  def newrepsonse(req)
    # Сохраняем слова вопроса для предотвращения повторения
    @last_question_words = Set.new(preprocess(req, true)) # true - сохраняем вопросительные слова
    
    if rand(0..10) == 5
      onlyforrefernce = "\n\n⚠️  Ответ был сгенерирован движком ответов NgramV1. Содержимое любых сообщений бота(за исключением /fetch и /amnesia) не являются информативными и генерируются исключительно для комедийного эффекта."
    else
      onlyforrefernce = ""
    end
    response = generate_response(req)
    return "#{response}" + onlyforrefernce
  end
  
  # Загрузка фраз из файла
  def load_vocabulary(filename = '../../../vocabulary.bot')
    if File.exist?(filename)
      @phrases = File.readlines(filename, encoding: 'utf-8')
                   .map(&:strip)
                   .reject { |line| line.empty? || line.start_with?('#') }
      
      puts "Загружено фраз: #{@phrases.size}"
      train_on_phrases
    else
      puts "Файл #{filename} не найден. Создан пустой файл."
      File.write(filename, "# Добавьте фразы для обучения бота\n# Каждая фраза на новой строке\n")
      @phrases = []
    end
  end
  
  # Обучение на всех фразах
  def train_on_phrases
    @phrases.each do |phrase|
      train(phrase, 2) # Больший вес для фраз из vocabulary
    end
    puts "Модель обучена на #{@phrases.size} фразах"
  end
  
  # Обучение модели на тексте с весом
  def train(text, weight = 1)
    words = preprocess(text, false) # при обучении используем обычную фильтрацию
    return if words.empty?

    # Считаем частоты слов
    words.each do |word|
      @word_freq[word] += weight
      @vocabulary.add(word)
    end
    
    # Считаем пары слов (биграммы)
    (0..words.size-2).each do |i|
      current_word = words[i]
      next_word = words[i+1]
      @word_pairs[current_word][next_word] += weight
    end
  end
  
  # Улучшенная предобработка текста с учетом типа
  def preprocess(text, keep_question_words = false)
    words = text.downcase
               .gsub(/[^\p{Cyrillic}\s\.\!\?\,]/, '') # сохраняем основные пунктуации
               .gsub(/\s+/, ' ')              # заменяем множественные пробелы на один
               .strip
               .split
               .reject { |word| word.empty? || word.size < 2 } # менее строгая фильтрация по длине
    
    if keep_question_words
      # Сохраняем вопросительные слова, но фильтруем остальные стоп-слова
      words.reject { |word| @stop_words.include?(word) && !@question_words.include?(word) }
    else
      # Обычная фильтрация
      words.reject { |word| @stop_words.include?(word) || @filler_words.include?(word) }
    end
  end
  
  # Генерация следующего слова с улучшенной логикой предотвращения повторений
  def next_word(current_word, used_words, recent_words = [])
    possibilities = @word_pairs[current_word]
    return smart_random_word_without_question(used_words, recent_words) if possibilities.empty?
    
    # Фильтруем уже использованные слова, короткие слова, недавно использованные и слова из вопроса
    available_words = possibilities.reject do |word, _| 
      used_words.include?(word) || 
      word.size < 3 || 
      recent_words.include?(word) ||
      question_word_repetition?(word)
    end
    
    if available_words.empty?
      # Если все возможные слова уже использованы, ослабляем фильтрацию (но все равно избегаем слов из вопроса)
      available_words = possibilities.reject { |word, _| word.size < 3 || question_word_repetition?(word) }
    end
    
    return smart_random_word_without_question(used_words, recent_words) if available_words.empty?
    
    total = available_words.values.sum
    return smart_random_word_without_question(used_words, recent_words) if total == 0
    
    rand_val = rand(total.to_f)
    cumulative = 0
    
    available_words.each do |word, count|
      cumulative += count
      return word if rand_val <= cumulative
    end
    
    smart_random_word_without_question(used_words, recent_words)
  end
  
  # Умный выбор случайного слова с учетом недавних слов и без слов из вопроса
  def smart_random_word_without_question(used_words = Set.new, recent_words = [])
    return "привет" if @vocabulary.empty?
    
    # Предпочитаем слова средней длины (4-8 букв), избегаем недавние и слова из вопроса
    good_words = @vocabulary.select do |word| 
      word.size.between?(4, 8) && 
      !used_words.include?(word) && 
      !@stop_words.include?(word) &&
      !recent_words.include?(word) &&
      !question_word_repetition?(word)
    end
    
    if good_words.any?
      good_words.sample
    else
      # Ослабляем условия, но все равно избегаем слов из вопроса
      (@vocabulary - used_words - @last_question_words)
        .reject { |word| @stop_words.include?(word) }
        .to_a.sample || 
      (@vocabulary - @last_question_words).to_a.sample
    end
  end
  
  # Генерация текста с улучшенной логикой предотвращения повторений
  def generate_text(start_word = nil, length = 6..10)
    return "Я еще не обучен. Добавьте фразы в vocabulary.bot" if @vocabulary.empty?
    
    current_word = start_word || find_relevant_start_without_question
    text_length = length.is_a?(Range) ? rand(length) : length
    text = [current_word]
    used_words = Set.new([current_word])
    recent_words = [current_word] # Трекер недавних слов
    
    (text_length - 1).times do
      current_word = next_word(current_word, used_words, recent_words.last(3))
      break if current_word.nil?
      
      # Проверяем, не повторяется ли слово слишком много раз
      if too_many_repetitions?(text, current_word) || question_word_repetition?(current_word)
        current_word = find_alternative_word(text, used_words, recent_words)
        break if current_word.nil?
      end
      
      text << current_word
      used_words.add(current_word)
      recent_words << current_word
      
      # Обновляем трекер недавних слов (ограничиваем размер)
      recent_words = recent_words.last(5) if recent_words.size > 5
      
      # Останавливаемся если достигли точки или есть законченная мысль
      break if should_stop_generation?(text, used_words)
    end
    
    # Постобработка для удаления повторений
    processed_text = remove_repetitions(text)
    processed_text.capitalize
  end
  
  # Проверка на повторение слов из вопроса
  def question_word_repetition?(word)
    @last_question_words.include?(word) && word.size > 3
  end
  
  # Проверка на слишком много повторений
  def too_many_repetitions?(text, new_word)
    return false if text.empty?
    
    # Проверяем, не повторяется ли слово слишком часто
    recent_count = text.last(4).count(new_word)
    recent_count >= 1 # Теперь запрещаем любые повторения в последних 4 словах
  end
  
  # Поиск альтернативного слова при повторениях
  def find_alternative_word(text, used_words, recent_words)
    return nil if text.empty?
    
    current_word = text.last
    possibilities = @word_pairs[current_word]
    
    # Ищем альтернативные слова, избегая повторений и слов из вопроса
    alternatives = possibilities.reject do |word, _|
      used_words.include?(word) || 
      recent_words.include?(word) || 
      word.size < 3 ||
      question_word_repetition?(word)
    end
    
    if alternatives.any?
      alternatives.keys.sample
    else
      # Если альтернатив нет, возвращаем умное случайное слово, избегая слов из вопроса
      smart_random_word_without_question(used_words, recent_words)
    end
  end
  
  # Удаление повторений из финального текста
  def remove_repetitions(words)
    return "" if words.empty?
    
    result = []
    previous_word = nil
    
    words.each do |word|
      if word == previous_word
        # Полностью пропускаем повторяющиеся слова
        next
      else
        result << word
        previous_word = word
      end
    end
    
    # Дополнительная проверка на общие повторения в тексте
    result = remove_global_repetitions(result)
    
    result.join(' ')
  end
  
  # Удаление глобальных повторений (слово встречается слишком часто во всем текста)
  def remove_global_repetitions(words)
    return words if words.size <= 3
    
    result = []
    seen_words = Set.new
    
    words.each do |word|
      # Если слово уже встречалось, пропускаем его
      if seen_words.include?(word)
        next
      else
        result << word
        seen_words.add(word)
      end
    end
    
    result
  end
  
  # Решение о остановке генерации
  def should_stop_generation?(text, used_words)
    return true if text.size >= 15 # Увеличили максимальную длину для вопросов
    return true if text.last.size < 3 && text.size > 4
    return true if nonsense_combination?(text)
    return true if too_many_repetitions?(text, text.last) # Останавливаемся при повторениях
    
    # Если много повторяющихся паттернов
    recent_words = text.last(4)
    recent_words.uniq.size < 2 && text.size > 5
  end
  
  # Проверка на бессмысленные комбинации
  def nonsense_combination?(words)
    return false if words.size < 3
    
    # Проверяем соотношение различных слов
    unique_ratio = words.uniq.size.to_f / words.size
    return true if unique_ratio < 0.4 && words.size > 5
    
    # Проверяем последние два слова
    last_two = words.last(2)
    last_two.all? { |w| w.size < 3 } && words.size > 3
  end
  
  # Поиск релевантного начального слова без слов из вопроса
  def find_relevant_start_without_question(input_words = [])
    return random_word_without_question if @phrases.empty?
    
    # Если есть входные слова, ищем наиболее релевантное (кроме слов из вопроса)
    if input_words.any?
      relevant_words = input_words.select do |word| 
        @vocabulary.include?(word) && word.size > 3 && !question_word_repetition?(word)
      end
      return relevant_words.sample if relevant_words.any?
    end
    
    # Или берем из частых осмысленных слов, исключая слова из вопроса
    frequent_words = @word_freq.sort_by { |_, count| -count }
                              .first(20)
                              .map(&:first)
                              .select do |word| 
                                word.size > 4 && 
                                !@stop_words.include?(word) &&
                                !question_word_repetition?(word)
                              end
    
    frequent_words.sample || random_word_without_question
  end
  
  # Случайное слово без слов из вопроса
  def random_word_without_question
    (@vocabulary - @last_question_words)
      .reject { |word| @stop_words.include?(word) }
      .to_a.sample || "привет"
  end
  
  # Определение типа запроса (вопрос/утверждение)
  def detect_question_type(input_words)
    return :question if input_words.any? { |word| @question_words.include?(word) }
    return :question if input_words.last&.include?('?')
    :statement
  end

  # Генерация ответа на фразу с улучшенным пониманием вопросов
  def generate_response(input_phrase, similarity_threshold = 0.3)
    return "Я еще не обучен. Добавьте фразы в vocabulary.bot" if @phrases.empty?
    
    input_words = preprocess(input_phrase, true) # сохраняем вопросительные слова
    return "Не понимаю ваш вопрос" if input_words.empty?
    
    # Определяем тип запроса
    question_type = detect_question_type(input_words)
    
    # Ищем наиболее релевантные фразы
    best_matches = find_best_matches(input_words, 5) # больше вариантов для вопросов
    
    if best_matches.any? { |match| match[:similarity] >= similarity_threshold }
      # Создаем временную модель для контекста
      relevant_phrases = best_matches.select { |m| m[:similarity] >= similarity_threshold }
                                    .map { |m| m[:phrase] }
      
      context_model = create_context_model(relevant_phrases)
      start_word = select_context_start_word(input_words, relevant_phrases, question_type)
      
      # Убеждаемся, что стартовое слово не из вопроса
      if question_word_repetition?(start_word)
        start_word = smart_random_word_without_question(Set.new, [])
      end
      
      response = context_model.generate_text(start_word, question_type == :question ? 6..12 : 4..8)
      
      # Для вопросов добавляем более осмысленное завершение
      if question_type == :question
        response = enhance_question_response(response, input_words)
      end
      
      response
    else
      # Общий ответ, избегая слов из вопроса
      start_word = find_relevant_start_without_question(input_words)
      response = generate_text(start_word, question_type == :question ? 5..10 : 4..6)
      
      if question_type == :question
        response = enhance_question_response(response, input_words)
      end
      
      response
    end
  end

  # Улучшение ответа на вопросы
  def enhance_question_response(response, input_words)
    # Анализируем тип вопроса и улучшаем ответ
    response_words = response.split
    
    if input_words.include?('что')
      # Для вопросов "что" стараемся дать определение
      if response_words.size > 2 && !response.downcase.start_with?('это')
        response = "Это " + response.downcase
      end
    elsif input_words.include?('кто')
      # Для вопросов "кто" 
      if response_words.size > 2 && !response.downcase.start_with?('это')
        response = "Это " + response.downcase
      end
    elsif input_words.include?('когда')
      # Для временных вопросов
      if !response.downcase.start_with?('обычно') && !response.downcase.start_with?('иногда')
        response = "Обычно " + response.downcase
      end
    elsif input_words.include?('где') || input_words.include?('куда')
      # Для пространственных вопросов
      if !response.downcase.start_with?('чаще') && !response.downcase.start_with?('обычно')
        response = "Чаще всего " + response.downcase
      end
    elsif input_words.include?('почему') || input_words.include?('зачем')
      # Для причинных вопросов
      if !response.downcase.start_with?('потому') && !response.downcase.start_with?('так')
        response = "Потому что " + response.downcase
      end
    elsif input_words.include?('как')
      # Для вопросов о способе
      if !response.downcase.start_with?('так') && !response.downcase.start_with?('обычно')
        response = "Обычно " + response.downcase
      end
    end
    
    response.capitalize
  end

  # Улучшенный поиск релевантных фраз с учетом вопросов
  def find_best_matches(input_words, limit = 3)
    matches = []
    question_type = detect_question_type(input_words)
    
    @phrases.each do |phrase|
      phrase_words = preprocess(phrase, false)
      similarity = calculate_similarity(input_words, phrase_words)
      
      # Бонус за соответствие типу (вопрос/ответ)
      if question_type == :question && (phrase.include?('?') || phrase_words.any? { |w| @question_words.include?(w) })
        similarity += 0.2
      end
      
      if similarity > 0.1
        matches << { phrase: phrase, similarity: similarity, words: phrase_words }
      end
    end
    
    matches.sort_by { |m| -m[:similarity] }.first(limit)
  end

  # Улучшенный выбор начального слова с учетом типа вопроса
  def select_context_start_word(input_words, relevant_phrases, question_type)
    # Собираем все слова из релевантных фраз
    all_words = relevant_phrases.flat_map { |phrase| preprocess(phrase, false) }
    
    # Для вопросов предпочитаем слова, которые часто отвечают на вопросы
    if question_type == :question
      question_answer_words = all_words.select do |word|
        word.size > 4 && !@stop_words.include?(word) && !question_word_repetition?(word)
      end
      
      return question_answer_words.sample if question_answer_words.any?
    end
    
    # Обычная логика выбора
    common_words = (input_words & all_words).select do |word| 
      word.size > 3 && !@stop_words.include?(word) && !question_word_repetition?(word)
    end
    
    return common_words.sample if common_words.any?
    
    informative_words = all_words.select do |word|
      word.size > 4 && @word_freq[word] > 1 && !@stop_words.include?(word) && !question_word_repetition?(word)
    end
    
    informative_words.sample || all_words.reject { |w| @stop_words.include?(w) || question_word_repetition?(w) }.sample
  end

  # Улучшенный расчет схожести с учетом вопросов
  def calculate_similarity(words1, words2)
    return 0 if words1.empty? || words2.empty?
    
    common_words = (words1 & words2)
    total_unique = (words1 | words2).size.to_f
    
    return 0 if total_unique == 0
    
    # Бонус за совпадение ключевых слов
    key_words1 = words1.select { |word| word.size > 3 }
    key_words2 = words2.select { |word| word.size > 3 }
    key_matches = (key_words1 & key_words2).size
    
    # Особый бонус за вопросительные слова
    question_matches = (words1 & @question_words.to_a & words2).size
    
    # Основная схожесть + бонусы
    (common_words.size / total_unique) + (key_matches * 0.1) + (question_matches * 0.3)
  end
  
  # Создание контекстной модели
  def create_context_model(phrases)
    model = self.class.new
    phrases.each { |phrase| model.train(phrase, 3) } # Больший вес для контекста
    model
  end
  
  # Добавление новой фразы в vocabulary
  def add_phrase(phrase, filename = 'vocabulary.bot')
    return
  end
  
  # Вывод статистики
  def stats
    {
      vocabulary_size: @vocabulary.size,
      total_words: @word_freq.values.sum,
      phrases_count: @phrases.size,
      most_common_words: @word_freq.sort_by { |_, count| -count }
                                  .first(10)
                                  .reject { |word, _| @stop_words.include?(word) },
      question_words: @question_words.to_a
    }
  end
end