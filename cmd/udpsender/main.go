package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	u, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069")
	if err != nil {
		fmt.Printf("cannot create sender: %v\n", err)
		return
	}

	con, err := net.DialUDP("udp", nil, u)
	if err != nil {
		fmt.Printf("unable to udp dial up: %v\n", err)
		return
	}
	defer con.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("unable to read input: %v\n", err)
		}

		con.Write([]byte(input))
	}
}
