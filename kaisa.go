package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const Version = "1.0.0"

func main() {
	fmt.Println("Kaisa", Version)
	host := "127.0.0.1"
	port := 3100

	log.Printf("Connecting to %v:%v", host, port)
	conn, err := net.Dial("tcp", "127.0.0.1:3100")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected successfully to ")
	}
	
	go perception(conn)

	stdinScanner := bufio.NewScanner(os.Stdin)
	var sExp string
	msgLen := make([]byte, 4)

	for {
		success := stdinScanner.Scan()
		if success {
			sExp = stdinScanner.Text()
		} else {
			conn.Close()
			log.Fatal(stdinScanner.Err())
		}

		if strings.ToLower(sExp) == "exit" {
			break
		}

		binary.BigEndian.PutUint32(msgLen, uint32(len(sExp)))
		conn.Write(msgLen)
		fmt.Fprint(conn, sExp)
	}

	conn.Close()
}

func perception(conn net.Conn) {
	reader := bufio.NewReader(conn)
	msgLenRaw := make([]byte, 4)
	
	for {
		for i := 0; i < 4; i++ {
			lenRaw, err := reader.ReadByte()
			fmt.Println(lenRaw)
			msgLenRaw[i] = lenRaw
			if err != nil {
				log.Print(err)
			}
		}

		msgLen := binary.BigEndian.Uint32(msgLenRaw)
		buf := make([]byte, msgLen)
		io.ReadFull(reader, buf)
		fmt.Println(string(buf))
	}
}
