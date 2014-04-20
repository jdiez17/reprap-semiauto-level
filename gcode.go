package main

import (
	"fmt"
	"sync"
)

var line int = 1

func envelopeCommand(command string) []byte {
	checksum := 0
	str := fmt.Sprintf("N%d %s ", line, command)
	line += 1

	for _, chr := range []byte(str) {
		checksum ^= int(chr)
	}
	checksum &= 0xff

	str = fmt.Sprintf("%s*%d\n", str, checksum)
	return []byte(str)
}

type GCodeCommandSender struct {
	Conn *Connection
	WG   *sync.WaitGroup
}

func NewGCodeCommandSender(c *Connection) *GCodeCommandSender {
	cs := &GCodeCommandSender{
		Conn: c,
		WG:   &sync.WaitGroup{},
	}

	ch := make(chan string)
	c.AddListener(ch)
	go cs.listenForAck(ch)

	return cs
}

func (g *GCodeCommandSender) listenForAck(ch <-chan string) {
	for msg := range ch {
		if msg == "ok\n" {
			g.WG.Add(-1)
		}
	}
}

func (g *GCodeCommandSender) SendCommand(command GCodeCommand) {
	str := command.ToString()
	if str[0] == 'G' && str != "G91" && str != "G90" { // buffered command: add 1 to wg
		g.WG.Add(1)
	}

	g.Conn.Outgoing <- []byte(envelopeCommand(str))
}
