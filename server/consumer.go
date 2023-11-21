package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Consumer interface {
	Subscribe()
	Consumer() chan []byte
}

type ConsumerOpts struct {
	Addr string
}

type GoatQConsumer struct {
	opts           ConsumerOpts
	MessageChannel chan []byte
}

func NewGoatQConsumer(opts ConsumerOpts, ch chan []byte) *GoatQConsumer {
	return &GoatQConsumer{
		opts:           opts,
		MessageChannel: ch,
	}
}

func (consumer *GoatQConsumer) Subscribe() {
	for {
		time.Sleep(time.Millisecond * 800)
		conn, err := net.Dial("tcp", consumer.opts.Addr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = conn.Write([]byte("READ "))
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go consumer.handleConn(conn)
	}
}

func (consumer *GoatQConsumer) handleConn(conn net.Conn) {
	buff := make([]byte, 1024)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println(err)
		}
		if err == io.EOF {
			fmt.Println("end of line reached")
			if err := conn.Close(); err != nil {
				fmt.Println(err)
			}
			return
		}
		consumer.handlePutToChannel(buff[:n])
		if err := conn.Close(); err != nil {
			fmt.Println(err)
		}
		return

		//conn.Close()
	}
}

func (consumer *GoatQConsumer) handlePutToChannel(msg []byte) bool {
	if strings.Split(string(msg), " ")[0] == "error" {
		//fmt.Printf("Error %s\n", string(msg))
		return false
	}
	consumer.MessageChannel <- msg
	return true
}

func (consumer *GoatQConsumer) Consume() chan []byte {
	return consumer.MessageChannel
}
