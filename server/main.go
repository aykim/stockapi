package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	counter     int
	counterLock sync.RWMutex
)

func main() {
	// Hello world, the web server

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("received hello request")
		io.WriteString(w, "Hello, world!\n")
	}

	countHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("received count request")
		// io.WriteString(w, fmt.Sprintf("Count is %d\n", counter))

		counterLock.RLock()
		defer counterLock.RUnlock()

		countString := "current Count: "
		if ctx.Err() != nil {
			countString = "final Count: "
		}

		io.WriteString(w, countString+strconv.FormatInt(int64(counter), 10))
	}

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/count", countHandler)

	go manageCount(ctx)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func manageCount(ctx context.Context) {
	timer := time.NewTimer(time.Second * 2)

	for {
		timer.Reset(time.Second * 2)

		select {
		case <-timer.C:
		case <-ctx.Done():
			// exits goroutine started above
			return
		}

		counterLock.Lock()

		if counter >= 200 {
			counter = 0
		} else {
			counter += 10
		}

		counterLock.Unlock()
	}
}
