package main

import (
	"fmt"
	"sync"
)

var (
	blocks [][]byte
	mu     sync.Mutex
)

func main() {
	const blockSize = 1000 * 1024 * 1024 // 1GB

	defer fmt.Print("Allocated: ", blockSize)

	for {
		go allocMem(blockSize)
	}

}

func allocMem(blocksize int) {
	mu.Lock()
	defer mu.Unlock()

	b := make([]byte, blocksize)
	blocks = append(blocks, b)

	select {}
}
