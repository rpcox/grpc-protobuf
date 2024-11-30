package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	var wg sync.WaitGroup
	done := make(chan interface{})

	wg.Add(1)
	go func() {
		i := 0
		for {
			select {
			case <-ticker.C:
				log.Printf("Tick %d\n", i)
				i++
				if i == 3 {
					wg.Done()
					return
				}
			case <-time.After(7 * time.Second):
				log.Println("Manual Tick")
			case _, ok := <- done:
				if !ok {
					log.Println("DONE")
					wg.Done()
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(75 * time.Second)
		defer ticker.Stop()

		select {
		case <-ticker.C:
			log.Println("closing done channel")
			close(done)
			wg.Done()
			return
		}
	}()

	wg.Wait()
	log.Println("exit main")

}
