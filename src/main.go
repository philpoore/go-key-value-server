package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

func error(msg string) {
	fmt.Println("Error :", msg)
}

var dataMap = make(map[string]string)
var dataMapMutex sync.Mutex

func get(key string) string {
	dataMapMutex.Lock()
	defer dataMapMutex.Unlock()

	return dataMap[key]
}

func deleteKey(key string) {
	dataMapMutex.Lock()
	defer dataMapMutex.Unlock()
	delete(dataMap, key)
}

func set(key string, value string) {
	dataMapMutex.Lock()
	defer dataMapMutex.Unlock()

	dataMap[key] = value
}

func handleConnection(c *net.TCPConn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	reader := bufio.NewReader(c)
	reply := func(msg string) {
		c.Write([]byte(msg + "\n"))
	}

	for {
		data, err := reader.ReadString('\n')

		if errors.Is(err, io.EOF) {
			error("Client disconnected")
			break
		}

		if err != nil {
			error("Reading data : " + err.Error())
			break
		}

		line := strings.TrimSpace(string(data))

		fmt.Println("line: ", line)

		parts := strings.Split(line, " ")
		cmd := strings.ToUpper(parts[0])

		if cmd == "GET" {
			key := parts[1]
			value := get(key)

			reply(value)
		} else if cmd == "DEL" {
			key := parts[1]
			deleteKey(key)

			reply("OK")
		} else if cmd == "SET" {
			key := parts[1]
			value := parts[2]

			set(key, value)

			reply("OK")
		} else if cmd == "QUIT" {
			break
		}
	}
	fmt.Printf("Closing %s\n", c.RemoteAddr().String())
	c.Close()
}

func main() {
	fmt.Println("Go Key-Value-Server")

	endpoint := net.TCPAddr{Port: 20000}
	listener, err := net.ListenTCP("tcp", &endpoint)

	if err != nil {
		error("Listen " + err.Error())
		return
	}

	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			error("Accept " + err.Error())
			return
		}

		go handleConnection(conn)
	}
}
