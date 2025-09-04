require 'set'
require 'socket'
require 'cgi'

require_relative "model.rb"

model = ChatLanguageModel.new
    
puts "Загрузка vocabulary.bot..."
model.load_vocabulary('vocabulary.bot')

server = TCPServer.new('localhost', 64000)
puts "NgramV1 service running: http://localhost:64000"


# Функция для ручного парсинга query string
def parse_query_string(query)
  params = {}
  return params if query.nil? || query.empty?
  
  query.split('&').each do |pair|
    key, value = pair.split('=', 2)
    if key && value
      # Декодируем URL-encoded значения
      key = CGI.unescape(key)
      value = CGI.unescape(value)
      params[key] = value
    end
  end
  params
end

loop do
  client = server.accept
  begin
    request = client.readpartial(2048)
    request_line = request.lines[0] || ""
    method, full_path, _http = request_line.split(' ')

    if method == 'GET'
      path, query = full_path.split('?', 2)
      if path == '/api/v1/ask'
        params = parse_query_string(query)
        value = params['param'] || ''
        body = model.newrepsonse(value)

        client.write "HTTP/1.1 200 OK\r\n"
        client.write "Content-Type: text/plain; charset=utf-8\r\n"
        client.write "Connection: close\r\n"
        client.write "Content-Length: #{body.bytesize}\r\n"
        client.write "\r\n"
        client.write body
      else
        body = "not found"
        client.write "HTTP/1.1 404 Not Found\r\n"
        client.write "Content-Type: text/plain; charset=utf-8\r\n"
        client.write "Connection: close\r\n"
        client.write "Content-Length: #{body.bytesize}\r\n"
        client.write "\r\n"
        client.write body
      end
    else
      body = "method not allowed"
      client.write "HTTP/1.1 405 Method Not Allowed\r\n"
      client.write "Content-Type: text/plain; charset=utf-8\r\n"
      client.write "Connection: close\r\n"
      client.write "Content-Length: #{body.bytesize}\r\n"
      client.write "\r\n"
      client.write body
    end
  rescue => e
    body = "Internal server error: #{e.message}"
    client.write "HTTP/1.1 500 Internal Server Error\r\n"
    client.write "Content-Type: text/plain; charset=utf-8\r\n"
    client.write "Connection: close\r\n"
    client.write "Content-Length: #{body.bytesize}\r\n"
    client.write "\r\n"
    client.write body
  ensure
    client.close
  end
end