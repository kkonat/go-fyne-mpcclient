package tcpclient

import (
	"bufio"
	"net"
	"sync"
	"time"
)

type Client struct {
	address       string //TCP serer address in the format "IP:port"
	singleRequest bool   // true for TCP servers which disconnect after a single request

	mtx           *sync.Mutex // lock for hw access
	conn          net.Conn    // (re-)established connection
	lastReconnect time.Time

	Online bool // server current online state
}

type Conf struct {
	Addr          string
	SingleRequest bool
}

func New(params Conf) *Client {
	clnt := &Client{
		address:       params.Addr,
		singleRequest: params.SingleRequest,
		mtx:           &sync.Mutex{}}

	if !clnt.singleRequest {
		clnt.reconnect()
	}
	return clnt
}

func (c *Client) reconnect() error {
	if !c.Online {
		if time.Since(c.lastReconnect) > time.Second {
			c.lastReconnect = time.Now()
			// log.Println("Reconnecting...")
			conn, err := net.DialTimeout("tcp", c.address, time.Millisecond*200)

			if err == nil {
				// log.Println(" OK connected...")
				c.conn = conn
				c.Online = true
			} else {
				c.Online = false
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
		if err != nil || !c.Online {
			return nil, err
		}
	}
	_, err := c.conn.Write([]byte(cmd + "\n"))
	if err != nil {
		if !c.singleRequest {
			c.conn.Close()
			c.Online = false
		}
		return nil, err
	}

	scanner := bufio.NewScanner(c.conn)
	resp, err := parseResponse(scanner)
	if err != nil {
		if !c.singleRequest {
			c.conn.Close()
			c.Online = false
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
