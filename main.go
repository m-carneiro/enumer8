package main

import (
	"enumer8/services"
	"flag"
	"fmt"
	"os"
)

func main() {
	domain := flag.String("domain", "", "Domain to check")
	wordlist := flag.String("wordlist", "", "Wordlist to use")

	flag.Parse()

	if *domain == "" || *wordlist == "" {
		fmt.Println("Please provide a domain AND a wordlist.")
		fmt.Println("Usage: enumer8 --domain=<domain> --wordlist=<wordlist>")
		os.Exit(1)
	}

	results := make(chan string)

	go func() {
		services.EnumerateSubdomains(*domain, *wordlist, results)
		close(results)
	}()

	for results := range results {
		fmt.Println(results)
	}
}
