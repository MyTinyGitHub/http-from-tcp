package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
  l, err := net.Listen("tcp", "127.0.0.1:42069")
  if err != nil {
    fmt.Println("unable to open file")
    os.Exit(201)
  }

   defer l.Close()


  for {
    a, err := l.Accept()

    if err != nil {
      fmt.Printf("error while accepting a connection: %v\n", err )
    }

    fmt.Println("connection was accepted")

    c := getLinesChannel(a)
    line := <-c
    fmt.Printf("%s\n", line)

    fmt.Println("connection was terminated")
  }
}

func getLinesChannel(r io.ReadCloser) <-chan string {
  lines := make(chan string)

  go func() {
    defer r.Close()
    defer close(lines)

    var read = make([]byte, 8)
    var line []byte

    for {
      readBytes, err := r.Read(read)
      if readBytes > 0 {
        for _, byt := range read {
          if byt == '\n' {
            lines <- string(line)
            line = make([]byte, 0)
            continue
          }
          line = append(line, byt)
        }
      }

      if err == io.EOF {
        lines <- string(line)
        break
      }
    }
  }()
  
  return lines
}
