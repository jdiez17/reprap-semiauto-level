package main

import (
	"bufio"
	"github.com/jdiez17/goserial"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type Connection struct {
	socket    io.ReadWriteCloser
	listeners []chan string
	Outgoing  chan []byte
}

func getSerialPorts() []string {
	result := []string{}
	items, _ := ioutil.ReadDir("/dev")

	for _, file := range items {
		if strings.Contains(file.Name(), "tty.") {
			result = append(result, "/dev/"+file.Name())
		}
	}

	return result
}

func NewConnection(port string, baud int) (*Connection, error) {
	c := &serial.Config{Name: port, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	listeners := make([]chan string, 1)
	outgoing := make(chan []byte)
	conn := &Connection{socket: s, listeners: listeners, Outgoing: outgoing}

	go conn.read()
	go conn.write()
	return conn, nil
}

func (c *Connection) read() {
	reader := bufio.NewReader(c.socket)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		for _, listener := range c.listeners {
			// nonblocking send
			select {
			case listener <- str:
			default:
			}
		}
	}
}

func (c *Connection) write() {
	for bytes := range c.Outgoing {
		_, err := c.socket.Write(bytes)
		log.Println(">>>", string(bytes))
		if err != nil {
			panic(err)
		}
	}
}

func (c *Connection) AddListener(ch chan string) {
	c.listeners = append(c.listeners, ch)
}
