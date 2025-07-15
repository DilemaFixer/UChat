package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DilemaFixer/UChat/src/chat"
)

func main() {
	switch os.Args[1] {
	case "-s":
		startAsServer()
	case "-c":
		startAsClient()
	default:
		fmt.Println("Usage: UChat server|client")
	}
}

func startAsServer() {
	chat := chat.NewUServer()
	chat.Start(":8080")
	go printRecivingMsg(chat)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		err := chat.Send(message)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func printRecivingMsg(c chat.UChat) {
	for {
		if !c.IsBusy() {
			continue
		}

		msg, err := c.Recive()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(msg)
	}
}

func startAsClient() {
	chat := chat.NewUClient()
	chat.Start("localhost:8080")
	go printRecivingMsg(chat)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		err := chat.Send(message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
