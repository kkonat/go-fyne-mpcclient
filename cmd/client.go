package main

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	serverAddr  string
	mtx         *sync.Mutex
	eventStream chan any
}

func NewClient(addr string) *Client {
	clnt := &Client{serverAddr: addr, mtx: &sync.Mutex{}, eventStream: make(chan any)}

	return clnt
}

func (s *Client) ask(cmd string) []string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	conn, err := net.DialTimeout("tcp", s.serverAddr, time.Millisecond*200)
	if err != nil {
		return nil // , errors.New("error connecting to: " + s.serverAddr)
	}
	defer conn.Close()

	fmt.Fprintln(conn, cmd) // sends command to the TCP server

	scanner := bufio.NewScanner(conn)
	resp := parseResponse(scanner)

	return resp
}

func parseResponse(scanner *bufio.Scanner) []string {
	resp := []string{}

	for scanner.Scan() {
		t := scanner.Text()
		if len(t) >= 3 && t[:3] == "ACK" {
			panic("unknown mpd server response: " + string(scanner.Text()))
		}
		if len(t) == 2 && t == "OK" {
			break
		} else {
			resp = append(resp, scanner.Text())
		}
	}
	return resp
}

func tryExtractString(data []string, key string, defaultVal string) string {
	for _, s := range data {
		if strings.HasPrefix(s, key) {
			return strings.Split(s, ": ")[1]
		}
	}
	return defaultVal // pass through
}

func tryExtractInt(data []string, key string, defaultVal int64) int64 {
	vStr := tryExtractString(data, key, "")
	if vStr != "" {
		value, err := strconv.ParseInt(vStr, 10, 64)
		if err == nil {
			return value
		}
	}
	return defaultVal
}
func tryExtractFloat(data []string, key string, defaultVal float64) float64 {
	vStr := tryExtractString(data, key, "")
	if vStr != "" {
		value, err := strconv.ParseFloat(vStr, 64)
		if err == nil {
			return value
		}
	}
	return defaultVal
}

func trkTimeToString(t float32) string {
	str := ""
	h := int(t) / 3600
	t -= float32(h * 3600)
	m := int(t) / 60
	t -= float32(m * 60)
	s := int(t)
	if h > 0 {
		str += fmt.Sprintf("%d:", h)
	}
	str += fmt.Sprintf("%d:%02d", m, s)
	return str
}

func calcHash(resp []string) uint32 {
	blob := strings.Join(resp, "")
	crc32q := crc32.MakeTable(0xD5828281)
	return crc32.Checksum([]byte(blob), crc32q)
}
