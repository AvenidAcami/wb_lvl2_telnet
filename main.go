package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"
)

// _, err := conn.Write([]byte(fmt.Sprintf("%s", str)))
func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	conn, err := net.Dial("tcp", "telehack.com:23")
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}

	defer conn.Close()

	for {
		select {
		case <-sigChan:
			return

		default:
			message := "NEWUSER"

			if _, err = conn.Write([]byte(message)); err != nil {
				fmt.Println(err)
				break
			}
			io.Copy(os.Stdout, conn)
			time.Sleep(1 * time.Second)
		}
	}

}
