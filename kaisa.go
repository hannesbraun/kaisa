package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const Version = "1.0.0"

func main() {
	fmt.Println("Kaisa", Version)

	hostFlag := flag.String("host", "localhost", "SimSpark server host address")
	portFlag := flag.Uint("port", 3100, "SimSpark server port")
	withPerception := flag.Bool("with-perception", false, "Print perception to console")
	flag.Parse()
	host := *hostFlag
	port := *portFlag

	simsparkAddress := host + ":" + strconv.Itoa(int(port))
	log.Println("Connecting to", simsparkAddress)
	conn, err := net.Dial("tcp", simsparkAddress)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected successfully to", simsparkAddress)
		defer conn.Close()
	}

	if *withPerception {
		go perception(conn)
	}

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
}

func perception(conn net.Conn) {
	reader := bufio.NewReader(conn)
	msgLenRaw := make([]byte, 4)

	for {
		for i := 0; i < 4; i++ {
			lenRaw, err := reader.ReadByte()
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
