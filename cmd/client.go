package main

import (
	"bufio"
	"net"
	"sync"
	"time"
)

type Client struct {
	serverAddr string
	mtx        *sync.Mutex
	conn       net.Conn
}

func NewClient(addr string) *Client {
	clnt := &Client{serverAddr: addr, mtx: &sync.Mutex{}}
	clnt.attemptReconnect()
	return clnt
}

func (s *Client) attemptReconnect() net.Conn {
	if s.conn == nil {
		conn, err := net.DialTimeout("tcp", s.serverAddr, time.Millisecond*200)
		if err == nil {
			s.conn = conn
		}
	}
	return s.conn
}

func (s *Client) Req(cmd string) []string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.attemptReconnect()
	if s.conn == nil {
		return nil // TODO: handle more detailed error connecting to server ?
	}

	_, err := s.conn.Write([]byte(cmd + "\n"))
	if err != nil {
		s.conn.Close()
		s.conn = nil
		return nil
	}

	scanner := bufio.NewScanner(s.conn)
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
