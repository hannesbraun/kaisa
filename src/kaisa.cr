require "option_parser"
require "socket"

module Kaisa
  VERSION = "0.1.0"

  server, port = "localhost", 3100

  option_parser = OptionParser.parse do |parser|
    parser.banner = "Kaisa - A manual Spark agent"

    parser.on "-v", "--version", "Show version" do
      puts "version #{VERSION}"
      exit
    end
    parser.on "-h", "--help", "Show help" do
      puts parser
      exit
    end
    parser.on "-s INSERVER", "--server INSERVER", "Server address" do |inserver|
      server = inserver
    end
    parser.on "-p INPORT", "--port=INPORT", "Server port" do |inport|
      port = inport
    end
    parser.missing_option do |option_flag|
      STDERR.puts "ERROR: #{option_flag} is missing something."
      STDERR.puts ""
      STDERR.puts parser
      exit(1)
    end
    parser.invalid_option do |option_flag|
      STDERR.puts "ERROR: #{option_flag} is not a valid option."
      STDERR.puts parser
      exit(1)
    end
  end

  puts "Connecting to #{server}:#{port}\n"

  socket = TCPSocket.new(server, port)
  # Effectors
  spawn do
    STDIN.blocking = true

    loop do
      msg = gets.not_nil!
      break if msg == "exit"

      if msg != nil
        socket << msg.size
        socket << msg
      end
    end

    socket.close
  end

  # Perception
  spawn do
    while !(socket.closed?)
      response = socket.gets
      puts response
    end
  end

  # Wait until socket is closed
  loop do
    Fiber.yield
    break if socket.closed?
  end
end
