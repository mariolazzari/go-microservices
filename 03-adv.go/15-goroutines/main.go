package main

import (
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

func processTruck(truck Truck) error {
	fmt.Printf("Start processing truck: %+v\n", truck)

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

func processFleet(trucks []Truck) error {
	var wg sync.WaitGroup

	for _, truck := range trucks {
		wg.Add(1)

		go func(t Truck) {
			if err := processTruck(truck); err != nil {
				log.Panicln(err)
			}
			wg.Done()
		}(truck)

	}

	wg.Wait()

	return nil
}

func main() {
	fleet := []Truck{
		&NormalTruck{id: "NT1", cargo: 0},
		&ElectricTruck{id: "ET1", cargo: 0, battery: 100},
		&NormalTruck{id: "NT2", cargo: 0},
		&ElectricTruck{id: "ET2", cargo: 0, battery: 100},
	}

	if err := processFleet(fleet); err != nil {
		fmt.Printf("Error processing fleet: %v\n", err)
		return
	}

	fmt.Println("All trucks processed")
}
