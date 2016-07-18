package main

import (
	"fmt"
	"net"
	"time"
)

const MAXPORTLEN uint16 = 65535

type DialAddr struct {
	ip string
	port uint16
}

func (d *DialAddr) buildAddr() string {
	return fmt.Sprintf("%s:%d",d.ip,d.port)
}

func (d DialAddr) check() bool {

	conn, err := net.DialTimeout("tcp",d.buildAddr(),time.Second * 5)
	if err != nil {
		return false
	}
	conn.Close()
	return true

}

func scan(ip string) (chan uint16, chan bool) {

	channel := make(chan uint16)
	done := make(chan bool)


	addr := DialAddr{}
	addr.ip = ip

	var i uint16

	for i = 0 ; i < MAXPORTLEN ; i++ {
		go func(port uint16) {
			addr.port = port
			if addr.check() {
				channel <- port
			}
			if port == MAXPORTLEN - 1 {
				time.Sleep(time.Second * 5)
				done <- true
			}
		}(i)
	}

	//time.Sleep(time.Second * 5)

	return channel, done
}



func main() {

	channel, done := scan("127.0.0.1")

	for {
		select {
		case p := <-channel:
			fmt.Println(p)
		case <-done:
				goto Done
		default :
			//fmt.Println("nothing yet")
		}
	}

	Done:
	fmt.Print("Done")

}