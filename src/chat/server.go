package chat

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const RejectConnectionMessage string = "REJECTED: Server is busy"

type UServer struct {
	addr      string
	listener  net.Listener
	partner   net.Conn
	mutex     sync.RWMutex
	errorChan chan error
	stopChan  chan struct{}
	isRunning bool
}

func NewUServer() *UServer {
	return &UServer{
		addr:      "",
		errorChan: make(chan error),
		stopChan:  make(chan struct{}),
	}
}

func (s *UServer) IsBusy() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.partner != nil
}

func (s *UServer) setPartner(conn net.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.partner = conn
}

func (s *UServer) getPartner() net.Conn {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.partner
}

func (s *UServer) getIsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isRunning
}

func (s *UServer) setIsRunning(isRunning bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.isRunning = isRunning
}

func (us *UServer) Start(addr string) error {
	if us.IsBusy() {
		return fmt.Errorf("Chan't start server , is alrady have connection")
	}

	if addr == "" || addr == " " {
		return fmt.Errorf("Try create UServer with empty or whitespace string")
	}

	us.addr = addr
	var err error
	us.listener, err = net.Listen("tcp", addr)

	if err != nil {
		return fmt.Errorf("Error when try create listener socket :", err)
	}

	us.isRunning = true
	go us.listenConnections()
	return nil
}

func (us *UServer) listenConnections() {
	defer us.listener.Close()

	for us.getIsRunning() {
		if tcpListener, ok := us.listener.(*net.TCPListener); ok {
			tcpListener.SetDeadline(time.Now().Add(1 * time.Second))
		}

		conn, err := us.listener.Accept()

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}

			select {
			case <-us.stopChan:
				return
			default:
				us.errorChan <- fmt.Errorf("Error accepting connection", err)
				continue
			}
		}

		if us.IsBusy() {
			rejectConnection(conn)
		}

		us.setPartner(conn)
	}
}

func rejectConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte(RejectConnectionMessage))
}

func (us *UServer) End() error {
	us.setIsRunning(false)
	close(us.stopChan)

	if us.listener != nil {
		us.listener.Close()
	}

	if partner := us.getPartner(); partner != nil {
		partner.Close()
	}

	return nil
}

func (us *UServer) Send(msg string) error {
	partner := us.getPartner()
	if partner == nil {
		return fmt.Errorf("Error geting partner to , when tru send msg :", msg)
	}

	_, err := partner.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("Error sending message to partner :", err)
	}
	return nil
}

func (us *UServer) Recive() (string, error) {
	partner := us.getPartner()
	if partner == nil {
		return "", fmt.Errorf("Error geting partner , when try recive msg")
	}

	buffer := make([]byte, 1024)
	n, err := partner.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("Error reading data from socket to buffer :", err)
	}
	return string(buffer[:n]), nil
}
