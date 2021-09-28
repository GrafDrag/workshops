//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
// Hint: time.Ticker can be used
// Hint 2: to calculate timediff for Advanced lvl use:
//
//  start := time.Now()
//	// your work
//	t := time.Now()
//	elapsed := t.Sub(start) // 1s or whatever time has passed

package main

import (
	"context"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	res := true
	ch := make(chan bool)
	var cancel context.CancelFunc
	ctx := context.Background()
	if !u.IsPremium {
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
	}

	go func() {
		process()
		ch <- true
	}()

	select {
	case <-ctx.Done():
		res = false
	case <-ch:
	}

	return res
}

func main() {
	RunMockServer()
}
