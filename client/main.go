package main

import (
	"fmt"
	"io"
	"net/http"
)

const helloPath string = "http://localhost:8080/hello"

func main() {

	var client http.Client
	resp, err := client.Get(helloPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer resp.Body.Close()

	result := make([]byte, 100)
	n, err := resp.Body.Read(result)
	if err != nil && err != io.EOF {
		fmt.Println("read error:", err.Error())
		return
	}

	if n > len(result) {
		fmt.Println("n too big")
		return
	}

	fmt.Printf("n is %d and result is %s\n", n, string(result))
}
