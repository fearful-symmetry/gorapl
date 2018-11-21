package gorapl

import (
	"fmt"
	"testing"
)

func TestPowerLimit(t *testing.T) {

	h, err := NewRAPL(0)
	if err != nil {
		t.Fatalf("Could not init: %s", err)
	}
	dat, err := h.ReadPowerLimit(DRAM)
	if err != nil {
		t.Fatalf("Could not read MSR: %s", err)
	}
	fmt.Printf("Got back: %#v\n", dat)

}
