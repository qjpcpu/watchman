package utils

import (
    "testing"
)

func TestWalk(t *testing.T) {
    list := Walk("/", 64)
    if len(list) == 0 {
        t.Fatal("It's not possible.")
    }
}
