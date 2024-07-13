package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Clock struct {
	mu        sync.Mutex
	ID        int
	timestamp int64
}

const Debug = false

func DPrintf(format string, a ...interface{}) {
	logger := log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds)
	if Debug {
		x := fmt.Sprintf(format, a...)
		logger.Println(x)
	}
}

func (c *Clock) tick(requestTimestamp int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.timestamp = max(c.timestamp, requestTimestamp) + 1
	DPrintf("process: %d, tick: %d", c.ID, c.timestamp)
	return c.timestamp
}
