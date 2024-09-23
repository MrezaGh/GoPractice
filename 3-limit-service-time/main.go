//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

const Threshold = 10

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

func (u *User) incrementTime() int64 {
	return atomic.AddInt64(&u.TimeUsed, 1)
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	done := make(chan int)
	go func() {
		process()
		close(done)
	}()

	for {
		select {
		case <-done:
			return true
		case <-time.Tick(1 * time.Second):
			if timeUsed := u.incrementTime(); timeUsed >= Threshold {
				return false
			}
		}
	}

}

func main() {
	RunMockServer()
}
