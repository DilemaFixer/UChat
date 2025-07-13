package main 

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main(){
	go startServer()
	go startClient()

	select{}
}

func startServer() {
	listen , err := net.Listen("tcp" , ":8080")
	if err != nil {
		fmt.Println("Error of starting server" , err)
	}

	for {
		conn , err := listen.Accept()
		if err != nil {
			fmt.Println("Have error when try accept connection to server :" , err)
			return
		}
		go handlingNewConnection(conn)
	}
}

func handlingNewConnection(connection net.Conn){
	defer connection.Close()

	go func() {
		buffer := make([]byte,1024)
		for {
			n, err := connection.Read(buffer)
			if err != nil {
				fmt.Println("Error when read bytes from socket :" , err)
				return
			}
			fmt.Println("Server get from user:" , string(buffer[:n]))
		}
	}()


	select{}
}

func startClient(){
	conn , err := net.Dial("tcp" , "localhost:8080")
	if err != nil {
		fmt.Println("Error when try start client :" , err)
		return
	}

	defer conn.Close()

	go func() {
        buffer := make([]byte, 1024)
        for {
            n, err := conn.Read(buffer)
            if err != nil { return }
            fmt.Printf("User get from server: %s", string(buffer[:n]))
        }
    }()

	 go func() {
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            message := scanner.Text() + "\n"
            conn.Write([]byte(message))
        }
    }()

	select {}
}


