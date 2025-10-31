package internal

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

func Connect(host, port string, timeout int) error {

	address := net.JoinHostPort(host, port)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к %s: %v", address, err)
	}
	defer conn.Close()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	done := make(chan struct{})

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "ошибка чтения из соединения: %v\n", err)
		}
		close(done)
	}()

	buf := make([]byte, 4096)

	for {
		select {
		case <-sigChan:
			fmt.Fprintln(os.Stderr, "\nОтключаюсь.")
			return nil
		case <-done:
			fmt.Fprintln(os.Stderr, "\nСоединение закрыто сервером.")
			return nil
		default:
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				if _, err := conn.Write(buf[:n]); err != nil {
					fmt.Fprintf(os.Stderr, "ошибка записи в соединение: %v", err)
					return err
				}
			}
			if err == io.EOF {
				fmt.Fprintln(os.Stderr, "\nЗавершение по Ctrl+D...")
				conn.Close()
				<-done
				return nil
			}
			if err != nil {
				return fmt.Errorf("ошибка чтения stdin: %v", err)
			}

		}
	}
}
