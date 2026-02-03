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
