package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
  file, err := os.Open("messages.txt")
  if err != nil {
    fmt.Println("unable to open file")
    os.Exit(201)
  }
  c := getLinesChannel(file)
  for {
    line := <- c
    fmt.Printf("read: %s\n", line)
  }
}

func getLinesChannel(r io.ReadCloser) <-chan string {
  lines := make(chan string)

  go func(r io.ReadCloser, c chan<- string) {
    var read = make([]byte, 8)
    var line []byte

    for {
      readBytes, err := r.Read(read)
      if readBytes > 0 {
        for _, byt := range read {
          if byt == '\n' {
            c <- string(line)
            line = make([]byte, 0)
            continue
          }
          line = append(line, byt)
        }
      }

      if err == io.EOF {
        os.Exit(0)
      }
    }
  } (r, lines)

  
  return lines
}
