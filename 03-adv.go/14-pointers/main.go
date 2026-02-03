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
