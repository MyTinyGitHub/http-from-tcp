package main

import (
	"fmt"
	"http-from-tcp/internal/request"
	"net"
	"os"
)

func main() {
  l, err := net.Listen("tcp", "127.0.0.1:42069")
  if err != nil {
    fmt.Println("unable to open port")
    os.Exit(201)
  }

   defer l.Close()

  fmt.Println("Starting to listen on port 42069")

  for {
    a, err := l.Accept()

    if err != nil {
      fmt.Printf("error while accepting a connection: %v\n", err )
    }

    fmt.Println("connection was accepted")

    //readAll(a)
    c, err := request.RequestFromReader(a)
    if err != nil {
      fmt.Printf("Erorr processing request: %v\n", err)
    } else {
      c.RequestLine.Print()
    }

    fmt.Println("connection was terminated")
  }
}
