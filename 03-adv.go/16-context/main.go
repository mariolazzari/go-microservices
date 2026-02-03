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

	for _, truck := range trucks {
		wg.Add(1)

		go func(t Truck) {
			if err := processTruck(ctx, truck); err != nil {
				log.Panicln(err)
			}
			wg.Done()
		}(truck)

	}

	wg.Wait()

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
}
