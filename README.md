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

### Testing

```go
package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("processTruck", func(t *testing.T) {
		t.Run("should load and reload truck cargo", func(t *testing.T) {
			nt := &NormalTruck{id: "nt-1"}
			et := &NormalTruck{id: "et-1"}

			err := processTruck(nt)
			if err != nil {
				t.Fatalf("Error processing truck: %s", err)
			}

			err = processTruck(et)
			if err != nil {
				t.Fatalf("Error processing truck: %s", err)
			}

			// asserting
			if nt.cargo != 0 {
				t.Fatalf("Normal truck cargo should be 0: %b", nt.cargo)
			}
			if et.cargo != 0 {
				t.Fatalf("Electric truck cargo should be 0: %b", et.cargo)
			}

		})
	})
}
```

### Pointers

```go
package main

import (
	"fmt"
	"log"
)

type NormalTruck struct {
	id    string
	cargo int
}

func fillTruck(nt *NormalTruck) {
	if nt == nil {
		fmt.Println("null truck pointer")
		return
	}
	nt.cargo = 100
}

func main() {
	id := 42
	log.Println("id", id)
	log.Println("&id", &id)

	ptrId := &id
	log.Println("ptrId", ptrId)
	log.Println("&ptrId", &ptrId)
	log.Println("*ptrId", *ptrId)

	// reset
	id = 0
	println("reset id")
	println("*ptrId", *ptrId)

	nt := NormalTruck{
		id:    "NT-0",
		cargo: 0,
	}

	// pass by reference
	fillTruck(&nt)
	fmt.Printf("NT: %+v\n", nt)

	var userId *int
	fmt.Println("Default pointer value:", userId)

	var nullTruck *NormalTruck
	fillTruck(nullTruck)
}
```

### Goroutines

```go
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
```

### Context

```go
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
```

### Concurrency with channels

```go
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
```

### Maps

```go
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
```

### Mutex

```go
package main

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"time"
)

type EvictCallback func(key string, value interface{})

type Options struct {
	withLogs      bool
	evictCallback EvictCallback
}

type LRU struct {
	keyMap map[string]*list.Element
	list   *list.List
	size   int
	opts   Options
	sync.RWMutex
}

type Item struct {
	key    string
	value  interface{}
	expiry time.Time
}

func NewLRUWithTTL(size int, opts Options) *LRU {
	if size <= 0 {
		return nil
	}

	lru := &LRU{
		size:   size,
		keyMap: make(map[string]*list.Element),
		list:   list.New(),
		opts:   opts,
	}

	return lru
}

// Set sets a new key and value into the cache with a ttl option.
// It returns true if a new element has been created.
func (l *LRU) Set(key string, value interface{}, ttl time.Duration) bool {
	l.Lock()
	defer l.Unlock()

	expiry := time.Now().Add(ttl)

	time.AfterFunc(ttl, func() {
		if l.opts.withLogs {
			log.Printf("Expiring: %s", key)
		}

		l.Lock()
		defer l.Unlock()

		if elem, ok := l.keyMap[key]; ok {
			item := elem.Value.(*Item)
			if time.Now().After(item.expiry) {
				l.list.MoveToBack(elem)

				if l.opts.withLogs {
					log.Printf("Elem moved back to list: %s", key)
				}
			}
		}
	})

	if elem, ok := l.keyMap[key]; ok {
		item := elem.Value.(*Item)
		item.value = value
		item.expiry = expiry

		l.list.MoveToFront(elem)

		if l.opts.withLogs {
			log.Printf("Elem %s updated to the front", key)
			l.printList()
		}

		return false
	}

	item := &Item{
		key:    key,
		value:  value,
		expiry: expiry,
	}

	elem := l.list.PushFront(item)

	l.keyMap[key] = elem

	if l.list.Len() > l.size {
		l.removeLastElement()
	}

	if l.opts.withLogs {
		log.Printf("Elem %s added to the front", key)
		l.printList()
	}

	return true
}

// Get gets by key and returns the value and if the entry is expired.
// If expired it is moved to the back of the list else it gets
// moved  to front as the most recently used
func (l *LRU) Get(key string) (value interface{}, expired bool) {
	l.Lock()
	defer l.Unlock()

	if elem, ok := l.keyMap[key]; ok {
		item := elem.Value.(*Item)

		expired := time.Now().After(item.expiry)
		if expired {
			l.list.MoveToBack(elem)
		} else {
			l.list.MoveToFront(elem)
		}

		return item.value, expired
	}

	return nil, false
}

func (l *LRU) printList() {
	for elem := l.list.Front(); elem != nil; elem = elem.Next() {
		item := elem.Value.(*Item)
		fmt.Printf("key: %s, value: %v \n", item.key, item.value)
	}
}

func (l *LRU) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	// remove from list
	l.list.Remove(e)

	// remove from map
	item := e.Value.(*Item)
	delete(l.keyMap, item.key)

	// evict element
	if l.opts.evictCallback != nil {
		l.opts.evictCallback(item.key, item.value)
	}
}

func (l *LRU) removeLastElement() {
	l.removeElement(l.list.Back())
}
```
