# Go microservices

## Design

### Monolith

[Monolith](https://martinfowler.com/bliki/StranglerFigApplication.html)

### Microservices

[Microservice](https://microservices.io/patterns/microservices.html)

### Communication styles

Pub/Sub

### Architecture

[Assemblage pattern](https://microservices.io/post/architecture/2023/02/09/assemblage-architecture-definition-process.html)

## Advanced Go

### Errors

```go
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
```

### Interfaces

```go
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
```