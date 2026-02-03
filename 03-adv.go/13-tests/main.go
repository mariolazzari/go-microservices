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

// oricess truck
func processTruck(t Truck) error {
	if err := t.LoadCargo(); err != nil {
		return fmt.Errorf("Error loading truck: %w", err)
	}

	if err := t.UnloadCargo(); err != nil {
		return fmt.Errorf("Error unloading truck: %w", err)
	}

	fmt.Printf("Progrssing track %+v\n", t)

	return nil
}

func main() {
	err := processTruck(&NormalTruck{id: "Normal-1"})
	if err != nil {
		log.Fatalf("Error processing truck: %s", err)
	}

	err = processTruck(&ElectricTruck{id: "Electric-1"})
	if err != nil {
		log.Fatalf("Error processing truck: %s", err)
	}
}
