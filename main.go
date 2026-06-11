package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// transport := &http.Transport{
	// 	ResponseHeaderTimeout: 5 * time.Second, // Time to wait for server&#39;s first response header
	// 	ExpectContinueTimeout: 1 * time.Second, // Time to wait for a response after sending an `Expect: 100-continue` header
	// }

	// client := &http.Client{
	// 	Transport: transport,
	// 	Timeout:   15 * time.Second, // Still advisable to set an overall timeout
	// }

	// resp, err := client.Get("https://api.example.com/data")
	// if err != nil {
	// 	fmt.Println("Request failed:", err)
	// 	return
	// }
	// defer resp.Body.Close()

	// fmt.Println("Response status:", resp.Status)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://api.github.com/")
	if err != nil {
		log.Fatalf("Request failed\n%v", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}
