package internal

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
)

func Connect(host, port string, timeout int) error {
	var message string

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	defer conn.Close()

	for {
		select {
		case <-sigChan:
			return nil

		default:

			scanner.Scan()

			message = scanner.Text()

			if _, err = conn.Write([]byte(message)); err != nil {
				fmt.Println(err)
				break
			}
			io.Copy(os.Stdout, conn)
		}
	}
}
