package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	m := make(map[string]int)
	m["a"] = 1

	a, ok := m["a"]
	if ok {
		fmt.Println("a:", a)
		delete(m, "a")
	}

	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(1 * time.Second)
			// race condition
			m[fmt.Sprintf("key-%d", i)] = i
		}(i)
	}
	wg.Wait()
	fmt.Println("Map:", m)
}
