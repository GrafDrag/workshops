//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, c chan *Tweet) {
	for tweet, err := stream.Next(); err == nil; tweet, err = stream.Next() {
		c <- tweet
	}

	close(c)
	wg.Done()
}

func consumer(c chan *Tweet) {
	for t := range c {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	start := time.Now()
	stream := GetMockStream()

	ch := make(chan *Tweet)
	wg.Add(2)

	// Producer
	go producer(stream, ch)

	// Consumer
	go consumer(ch)

	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
