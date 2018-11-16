package gorapl

import (
	"fmt"
	"testing"
)

func TestPowerLimit(t *testing.T) {

	ret, err := ReadPowerLimit(DRAM)
	if err != nil {
		t.Fatalf("Could not read MSR: %s", err)
	}

	fmt.Printf("Got back: %#v\n", ret)

}
