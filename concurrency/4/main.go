package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const stringToSearch = "concurrency"

var sites = []string{
	"https://google.com",
	"https://itc.ua/",
	"https://twitter.com/concurrencyinc",
	"https://twitter.com/",
	"http://localhost:8000",
	"https://github.com/bradtraversy/go_restapi/blob/master/main.go",
	"https://www.youtube.com/",
	"https://postman-echo.com/get",
	"https://en.wikipedia.org/wiki/Concurrency_(computer_science)#:~:text=In%20computer%20science%2C%20concurrency%20is,without%20affecting%20the%20final%20outcome.",
}

type SiteData struct {
	data []byte
	uri  string
}

var wg = sync.WaitGroup{}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	resultsCh := make(chan SiteData, len(sites))
	wg.Add(len(sites))

	for _, site := range sites {
		go Worker(ctx, site, resultsCh)
	}

	for ch := range resultsCh {
		if strings.Contains(string(ch.data), stringToSearch) {
			fmt.Printf("'%s' string is found in %s\n", stringToSearch, ch.uri)
			cancel()
			break
		} else {
			fmt.Printf("Nothing found in %s\n", ch.uri)
		}
	}

	fmt.Println("exiting from searcher...")

	wg.Wait()
	// give one second to validate if all other goroutines are closed
	//time.Sleep(time.Second)
}

func Worker(ctx context.Context, uri string, ch chan SiteData) {
	defer wg.Done()
	done := make(chan bool)
	sd := SiteData{
		uri: uri,
	}
	go func() {
		fmt.Printf("starting sending request to %s\n", sd.uri)
		sd.data = Reader(sd.uri)
		done <- true
	}()

	select {
	case <-done:
		ch <- sd
	case <-ctx.Done():
		fmt.Printf("Get \"%s\": context canceled\n", sd.uri)
	}
}

func Reader(uri string) []byte {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return bodyBytes
}
