package main

import (
	"encoding/json"
	"log"
	"math"
	"net"
	"os"
	"time"
)

const (
	SAMPLES_SIZE = 10000
	NODES_FILE   = "nodes.txt"
)

func FindShortestPath(ori string, dest string, encoder *json.Encoder, decoder *json.Decoder) (*ResponsePayload, *time.Duration, error) {
	requestPayload := RequestPayload{
		Ori:  ori,
		Dest: dest,
	}
	start := time.Now()
	if err := encoder.Encode(requestPayload); err != nil {
		return nil, nil, err

	}
	rtt := time.Since(start)

	responsePayload := ResponsePayload{}
	if err := decoder.Decode(&responsePayload); err != nil {
		return nil, nil, err
	}

	return &responsePayload, &rtt, nil
}

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal("undefined HOST")
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

	ori := "A"
	dest := "E"
	var samples []time.Duration
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	for i := 0; i < SAMPLES_SIZE; i++ {
		log.Printf("sending request to find the shortest path between %s and %s", ori, dest)
		res, rtt, err := FindShortestPath(ori, dest, encoder, decoder)
		if err != nil {
			log.Fatal(err)
		}
		samples = append(samples, *rtt)
		log.Printf("shortest path received %v", res.Path)
	}

	var mean float64
	for _, sample := range samples {
		mean += float64(sample)
	}
	mean = mean / float64(len(samples))

	var sd float64
	for _, sample := range samples {
		sd += math.Pow((float64(sample) - mean), 2)
	}
	sd = math.Sqrt(sd / float64(len(samples)))

	log.Printf("average RTT is %.2f (+- %.2f)", mean, sd)
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
