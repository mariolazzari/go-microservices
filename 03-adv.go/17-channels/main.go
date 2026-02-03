package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Truck interface {
	LoadCargo() error
	UnloadCargo() error
}

type NormalTruck struct {
	id    string
	cargo int
}

type ElectricTruck struct {
	id      string
	cargo   int
	battery float64
}

func (nt *NormalTruck) LoadCargo() error {
	nt.cargo++
	return nil
}

func (nt *NormalTruck) UnloadCargo() error {
	nt.cargo--
	return nil
}

func (et *ElectricTruck) UnloadCargo() error {
	et.cargo++
	et.battery++
	return nil
}

func (et *ElectricTruck) LoadCargo() error {
	et.cargo--
	et.battery--
	return nil
}

func processTruck(ctx context.Context, truck Truck) error {
	fmt.Printf("Start processing truck: %+v\n", truck)

	// access user id
	userId := ctx.Value("userId")
	log.Println("User ID:", userId)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	delay := 3 * time.Second
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		break
	}

	// emulate long operation
	time.Sleep(1 * time.Second)

	err := truck.LoadCargo()
	if err != nil {
		return fmt.Errorf("error loading cargo: %w", err)
	}

	err = truck.UnloadCargo()
	if err != nil {
		return fmt.Errorf("error unloading cargo: %w", err)
	}

	fmt.Printf("Finish processing truck: %+v\n", truck)

	return nil
}

func processFleet(ctx context.Context, trucks []Truck) error {
	var wg sync.WaitGroup
	errorsChan := make(chan error, len(trucks))

	for _, truck := range trucks {
		wg.Add(1)

		go func(t Truck) {
			if err := processTruck(ctx, truck); err != nil {
				log.Panicln(err)
				errorsChan <- err
			}
			wg.Done()
		}(truck)

	}

	wg.Wait()

	close(errorsChan)

	// select {
	// case err := <-errorsChan:
	// 	return err
	// default:
	// 	return nil
	// }

	var errs []error
	for err := range errorsChan {
		log.Printf("Error processing truck: %v\n", err)
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("fleet processing has %d errors", len(errs))
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "userId", 42)

	fleet := []Truck{
		&NormalTruck{id: "NT1", cargo: 0},
		&ElectricTruck{id: "ET1", cargo: 0, battery: 100},
		&NormalTruck{id: "NT2", cargo: 0},
		&ElectricTruck{id: "ET2", cargo: 0, battery: 100},
	}

	if err := processFleet(ctx, fleet); err != nil {
		fmt.Printf("Error processing fleet: %v\n", err)
		return
	}

	fmt.Println("All trucks processed")

	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(1 * time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
