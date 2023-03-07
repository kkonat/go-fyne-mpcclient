package tcpclient

import (
	"bufio"
	"log"
	"strings"
	"testing"
)

const connected = `OK MPD 0.16.0`
const currentsong = `file: Playlisty/SpotifyPotW/23/02.01/Shannon & The Clams - Do I Wanna Stay.mp3
Last-Modified: 2023-02-07T14:34:34Z
Time: 269
Artist: Shannon & The Clams
AlbumArtist: Shannon & The Clams
Title: Do I Wanna Stay
Album: Year Of The Spider
Track: 1/13
Date: 2021-08-20
Genre: bay area indie
Disc: 1/1
Pos: 29
Id: 29
OK`

func TestScanner(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(currentsong))
	resp, err := parseResponse(scanner)
	if err != nil {
		t.Error("error parsing response")
	}
	l12 := strings.Split(currentsong, "\n")[12]
	if resp[12] != l12 {
		t.Log("\n" + strings.Join(resp, "\n"))
		t.Error("response not matching input data")
	}
}

func TestSendCtrlCmd(t *testing.T) {
	mpc := New(Conf{Addr: "192.168.0.95:6600", SingleRequest: false})

	resp, err := mpc.Request("currentsong")
	if err != nil {
		t.Error("error sending request")
	}
	log.Println(resp)

	if resp == nil {
		t.Error("Test fail: nil response")
	}
	if resp[len(resp)-1][:3] != "Id:" {
		t.Error("Last item should start with `Id:`")
	}
}
