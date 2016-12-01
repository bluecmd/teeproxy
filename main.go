package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	listen = flag.String("listen", "localhost:1234", "Address to listen on")
	addr_a = flag.String("a", "localhost:1235", "First address to send traffic to")
	addr_b = flag.String("b", "localhost:1236", "Second address to send traffic to")
)

func sender(host string, dc chan []byte) {
	c, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Printf("Error dialing %s: %s\n", host, err.Error())
		return
	}
	defer c.Close()
	for {
		d := <-dc
		if len(d) == 0 {
			return
		}
		c.Write(d)
	}
}

func handle(c net.Conn) {
	buf := make([]byte, 1024)
	defer c.Close()

	a_chan := make(chan []byte)
	b_chan := make(chan []byte)

	go sender(*addr_a, a_chan)
	go sender(*addr_b, b_chan)

	for {
		l, err := c.Read(buf)
		if l <= 0 && err != nil {
			a_chan <- []byte{}
			b_chan <- []byte{}
			return
		}

		a_chan <- buf[:l]
		b_chan <- buf[:l]
	}
}

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", *listen)
	if err != nil {
		panic(err);
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}

		go handle(c)
	}
}
