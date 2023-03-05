package main

import (
	"bufio"
	"net"
	"sync"
	"time"
)

type TCPClient struct {
	address       string      //TCP serer address in the format "IP:port"
	mtx           *sync.Mutex // lock for hw access
	conn          net.Conn    // (re-)established connection
	singleRequest bool        // true for TCP servers which disconnect after a single request
	online        bool        // server current online state
	lastReconnect time.Time
}

type TCPClientParms struct {
	addr          string
	singleRequest bool
}

func NewClient(params TCPClientParms) *TCPClient {
	clnt := &TCPClient{address: params.addr, mtx: &sync.Mutex{}, singleRequest: params.singleRequest}
	if !clnt.singleRequest {
		clnt.reconnect()
	}
	return clnt
}

func (c *TCPClient) reconnect() error {
	if !c.online {
		if time.Since(c.lastReconnect) > time.Second {
			c.lastReconnect = time.Now()
			// log.Println("Reconnecting...")
			conn, err := net.DialTimeout("tcp", c.address, time.Millisecond*200)

			if err == nil {
				// log.Println(" OK connected...")
				c.conn = conn
				c.online = true
			} else {
				c.online = false
				// log.Println("Failed to connect to server :", err.Error())
				return err
			}
		}
	}
	return nil
}

func (c *TCPClient) Request(cmd string) ([]string, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.singleRequest {
		conn, err := net.DialTimeout("tcp", c.address, time.Millisecond*200)
		if err == nil {
			// log.Println(" OK connected...")
			c.conn = conn
		} else {
			return nil, err
		}
	} else {
		err := c.reconnect()
		if err != nil || !c.online {
			return nil, err
		}
	}
	_, err := c.conn.Write([]byte(cmd + "\n"))
	if err != nil {
		if !c.singleRequest {
			c.conn.Close()
			c.online = false
		}
		return nil, err
	}

	scanner := bufio.NewScanner(c.conn)
	resp, err := parseResponse(scanner)
	if err != nil {
		if !c.singleRequest {
			c.conn.Close()
			c.online = false
		}
		return nil, err
	}
	if c.singleRequest {
		c.conn.Close()
	}
	return resp, nil
}

func parseResponse(scanner *bufio.Scanner) ([]string, error) {
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
	err := scanner.Err()
	return resp, err
}
