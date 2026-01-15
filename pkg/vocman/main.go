package vocman

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetChatVocabulary(chatID string, mainVocabulary []string) []string {
	// Проверяем, находится ли чат в ignore-list
	if isChatNeglected("neglected.bot", chatID) {
		// Читаем локальный словарь для этого чата
		localDictFile := getLocalDictFilename(chatID)
		return readVocabularyFromFile(localDictFile)
	}

	// Если чат не в ignore-list, используем основной словарь
	return mainVocabulary
}

func readVocabularyFromFile(filename string) []string {
	var vocabulary []string

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Файл словаря %s не найден", filename)
		} else {
			log.Printf("Ошибка открытия файла %s: %v", filename, err)
		}
		return vocabulary
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			vocabulary = append(vocabulary, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения файла %s: %v", filename, err)
	}

	log.Printf("Загружено %d фраз из %s", len(vocabulary), filename)
	return vocabulary
}

func AppendMessageToFile(filename string, message string, chatID string) {
	log.Printf("check ChatID: %s", chatID)

	if strings.Contains(message, "http://") ||
		strings.Contains(message, "https://") ||
		strings.Contains(message, "@") {
		log.Println("Links found. Skip")
		return
	}

	if isChatNeglected("neglected.bot", chatID) {
		log.Printf("Chat: %s в ignore-list. Создаем локальный словарь", chatID)

		// Создаем имя файла для локального словаря на основе chatID
		localDictName := getLocalDictFilename(chatID)
		appendToLocalDictionary(localDictName, message)
		return
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Cannot open file %s: %v", filename, err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(message + "\n"); err != nil {
		log.Printf("Write error %s: %v", filename, err)
	} else {
		log.Printf("Message write in %s", filename)
	}
}

func isChatNeglected(filename, targetID string) bool {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("neglected.bot not found")
		}
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == targetID {
			log.Printf("Chat %s найден в ignore-list", targetID)
			return true
		}
	}
	return false
}

func getLocalDictFilename(chatID string) string {
	// Создаем папку для локальных словарей, если ее нет
	localDictsDir := "local_dicts"
	if err := os.MkdirAll(localDictsDir, 0755); err != nil {
		log.Printf("Ошибка создания папки %s: %v", localDictsDir, err)
		return "local_dict.txt" // fallback
	}

	// Создаем имя файла на основе chatID (убираем специальные символы)
	safeChatID := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			return r
		}
		return '_'
	}, chatID)

	return filepath.Join(localDictsDir, safeChatID+".txt")
}

func appendToLocalDictionary(filename string, message string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Cannot open local dictionary %s: %v", filename, err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(message + "\n"); err != nil {
		log.Printf("Write error in local dictionary %s: %v", filename, err)
	} else {
		log.Printf("Message written to local dictionary: %s", filename)
	}
}
