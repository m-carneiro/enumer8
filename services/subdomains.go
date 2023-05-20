package services

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

const concurrency = 100

func checkSubdomain(subdomain string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := net.LookupHost(subdomain)
	if err == nil {
		results <- subdomain
	}
}

func EnumerateSubdomains(domain string, wordlist string, results chan<- string) {
	file, err := os.Open(wordlist)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(file)

	scanner := bufio.NewScanner(file)

	var waitGroup sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for scanner.Scan() {
		waitGroup.Add(1)
		semaphore <- struct{}{}

		go func(subdomain string) {
			defer func() { <-semaphore }()
			checkSubdomain(subdomain, results, &waitGroup)
		}(scanner.Text() + "." + domain)
	}

	for scanner.Scan() {
		subdomain := scanner.Text() + "." + domain
		_, err := net.LookupHost(subdomain)
		if err != nil {
			fmt.Println(subdomain, "not found")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}
