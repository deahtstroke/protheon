package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"time"
)

type Job struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

func runServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server listening on :%s\n", port)

	var clients []net.Conn

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Accept error: ", err)
				continue
			}
			log.Printf("Worker connected: %s\n", conn.RemoteAddr())
			clients = append(clients, conn)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	for t := range ticker.C {
		job := Job{
			ID:   rand.Int(),
			Data: fmt.Sprintf("Job created at %s", t.Format(time.RFC3339)),
		}
		msg, _ := json.Marshal(job)
		for _, c := range clients {
			_, err := c.Write(append(msg, '\n'))
			if err != nil {
				log.Printf("Failed to send job to %s: %v", c.RemoteAddr(), err)
			}
		}
		log.Printf("Emmited job: %+v\n", job)
	}
}

func runWorker(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	log.Printf("Connected to server at %s\n", addr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var job Job
		if err := json.Unmarshal(scanner.Bytes(), &job); err != nil {
			log.Println("error decoding job: ", err)
			continue
		}
		fmt.Printf("Worker got job: %+v\n", job)
	}
}

func main() {
	role := flag.String("role", "worker", "server or worker")
	port := flag.String("port", "9000", "server port")
	addr := flag.String("addr", "127.0.0.1:9000", "server address for worker")
	flag.Parse()

	if *role == "server" {
		runServer(*port)
	} else if *role == "worker" {
		runWorker(*addr)
	} else {
		fmt.Println("Unknown role, use --role=server or --role=worker")
		os.Exit(1)
	}
}
