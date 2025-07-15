package chat

import (
	"fmt"
	"net"
)

type UClient struct {
	server net.Conn
}

func NewUClient() *UClient {
	return &UClient{}
}

func (u *UClient) IsBusy() bool {
	return u.server != nil
}

func (u *UClient) Start(addr string) error {
	if u.IsBusy() {
		return fmt.Errorf("UClient is busy")
	}

	if addr == "" || addr == " " {
		return fmt.Errorf("Invalid address")
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("Get errror , when try start UClient :", err)
	}
	u.server = conn
	return nil
}

func (c *UClient) Send(message string) error {
	if !c.IsBusy() {
		return fmt.Errorf("UClient is not connected")
	}

	if message == "" || message == " " {
		return fmt.Errorf("Invalid message")
	}

	_, err := c.server.Write([]byte(message + "\n"))
	if err != nil {
		return fmt.Errorf("Get error when try send message to server", err)
	}

	return nil
}

func (u *UClient) End() error {
	if u.server == nil {
		return fmt.Errorf("UClient is not connected")
	}

	err := u.server.Close()
	if err != nil {
		return fmt.Errorf("Get error when try end UClient", err)
	}
	u.server = nil
	return nil
}

func (u *UClient) Recive() (string, error) {
	if u.server == nil {
		return "", fmt.Errorf("UClient is not connected")
	}

	buffer := make([]byte, 1024)
	n, err := u.server.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("Get error when try recive message from server", err)
	}

	return string(buffer[:n]), nil
}
