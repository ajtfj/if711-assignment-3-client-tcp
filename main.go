package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	SAMPLES_SIZE = 10000
	NODES_FILE   = "nodes.txt"
)

var (
	nodes []string
)

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal("undefined PORT")
	}

	if err := setup(); err != nil {
		log.Fatal(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer closeTCPConnection(conn)

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	for i := 0; i < SAMPLES_SIZE; i++ {
		rand.Seed(time.Now().UnixNano())
		ori := nodes[rand.Intn(len(nodes))]
		dest := nodes[rand.Intn(len(nodes))]

		requestPayload := RequestPayload{
			Ori:  ori,
			Dest: dest,
		}
		if err := encoder.Encode(requestPayload); err != nil {
			log.Fatal(err)
		}

		responsePayload := ResponsePayload{}
		if err := decoder.Decode(&responsePayload); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Shortest path between %s and %s is: %v\n", ori, dest, responsePayload.Path)
	}
}

func setup() error {
	file, err := os.ReadFile(NODES_FILE)
	if err != nil {
		return err
	}

	nodes = strings.Split(string(file), " ")

	return nil
}

type RequestPayload struct {
	Ori  string `json:"ori"`
	Dest string `json:"dest"`
}

type ResponsePayload struct {
	Path         []string      `json:"path"`
	CalcDuration time.Duration `json:"calc-duration"`
}

type ResponseErrorPayload struct {
	Message string `json:"message"`
}

func closeTCPConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
