package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

var (
	host        = flag.String("host", "localhost", "host to connect to")
	port        = flag.Int("port", 80, "port to connect to")
	num         = flag.Int("n", 1, "number of times to request and drop")
	concurrency = flag.Int("c", 1, "how many threads to send requests on")
)

func main() {
	flag.Parse()

	if *concurrency > *num {
		fmt.Println("error: c cannot be greater than n")
		os.Exit(1)
	}

	msgs := make(chan struct{}, 0)
	quit := make(chan struct{}, 0)
	wg := new(sync.WaitGroup)
	// Create the goroutine pool.
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-msgs:
					if err := sendAndDrop(); err != nil {
						log.Fatal(err)
					}
				case <-quit:
					wg.Done()
					return
				}
			}
		}()
	}

	for i := 0; i < *num; i++ {
		msgs <- struct{}{}
	}
	close(quit)
	wg.Wait()

}

func sendAndDrop() error {
	conn, err := net.Dial("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		return err
	}
	fmt.Println("connected")
	_, err = fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	if err != nil {
		return err
	}
	fmt.Println("sent")
	return conn.Close()
}
