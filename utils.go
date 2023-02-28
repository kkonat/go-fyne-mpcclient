package main

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"net"
	"os"
	"strconv"
	"strings"
)

func sendCtrlCmd(server, cmd string) []string {
	conn, err := net.Dial("tcp", server)

	var resp []string

	if err != nil {
		fmt.Println("Error connecting to host", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Fprintln(conn, cmd)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {

		if scanner.Text() == "OK" {
			break
		} else {
			resp = append(resp, scanner.Text())
		}
	}
	return resp
}
func extract(data []string, pattern string) string {
	for _, s := range data {
		if strings.HasPrefix(s, pattern) {
			return strings.Split(s, ": ")[1]
		}
	}

	return ""
}
func getVolume() (int64, error) {

	resp := sendCtrlCmd(mpdSrv, "status")
	v := extract(resp, "volume")

	return strconv.ParseInt(v, 10, 32)
}
func getTrackLen() (int64, error) {

	resp := sendCtrlCmd(mpdSrv, "currentsong")
	v := extract(resp, "Time :")

	return strconv.ParseInt(v, 10, 32)
}

func getTrackDataHash() uint32 {
	resp := sendCtrlCmd(mpdSrv, "currentsong")
	album := extract(resp, "Album:")
	artist := extract(resp, "Artist")
	track := extract(resp, "Title")
	blob := album + artist + track
	crc32q := crc32.MakeTable(0xD5828281)
	return crc32.Checksum([]byte(blob), crc32q)
}
