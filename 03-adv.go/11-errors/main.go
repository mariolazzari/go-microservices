package main

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrTruckNotFound  = errors.New("truck not found")
)

type Truck struct {
	id string
}

func (t *Truck) LoadCargo() error {
	return ErrTruckNotFound
}

func (t *Truck) UnloadCargo() error {
	return ErrTruckNotFound
}

// oricess truck
func processTruck(t Truck) error {
	fmt.Printf("Processing truck %s\n", t.id)

	if err := t.LoadCargo(); err != nil {
		return fmt.Errorf("Error loading truck: %w", err)
	}

	if err := t.UnloadCargo(); err != nil {
		return fmt.Errorf("Error unloading truck: %w", err)
	}

	return ErrNotImplemented
}

func main() {
	trucks := []Truck{
		{id: "Truck-1"},
		{id: "Truck-2"},
		{id: "Truck-3"},
	}

	for _, truck := range trucks {
		fmt.Printf("Truck %s arrived\n", truck.id)

		// err := processTruck(truck)
		// if err != nil {
		// 	log.Fatalf("Error processing truck: %s", err)
		// }

		if err := processTruck(truck); err != nil {
			log.Fatalf("Error processing truck: %s", err)
		}

	}

}
