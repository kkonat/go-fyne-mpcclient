package main

import (
	"bufio"
	"net"
	"sync"
	"time"
)

type Client struct {
	serverAddr    string
	mtx           *sync.Mutex
	conn          net.Conn
	connectOnce   bool
	online        bool
	lastReconnect time.Time
}

func NewClient(addr string, connectOnce bool) *Client {
	clnt := &Client{serverAddr: addr, mtx: &sync.Mutex{}, connectOnce: connectOnce}
	if !connectOnce {
		clnt.reconnect()
	}
	return clnt
}

func (c *Client) reconnect() error {
	if !c.online {
		if time.Since(c.lastReconnect) > time.Second {
			c.lastReconnect = time.Now()
			// log.Println("Reconnecting...")

			conn, err := net.DialTimeout("tcp", c.serverAddr, time.Millisecond*200)

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

func (c *Client) Request(cmd string) ([]string, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.connectOnce {
		conn, err := net.DialTimeout("tcp", c.serverAddr, time.Millisecond*200)
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
		if !c.connectOnce {
			c.conn.Close()
			c.online = false
		}
		return nil, err
	}

	scanner := bufio.NewScanner(c.conn)
	resp, err := parseResponse(scanner)
	if err != nil {
		if !c.connectOnce {
			c.conn.Close()
			c.online = false
		}
		return nil, err
	}
	if c.connectOnce {
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
