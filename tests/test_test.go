package test

import (
	"testing"
)

func TestSum (t *testing.T) {
    n := 1+1
    if n != 2 {
        t.Errorf("result should be 2")
    }
}

func TestMin (t *testing.T) {
    n := 1-1
    if n != 0 {
        t.Errorf("result should be 0")
    }
}
