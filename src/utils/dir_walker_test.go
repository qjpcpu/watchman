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
func TestCpu(t *testing.T) {
    _, err := Cpu()
    if err != nil {
        t.Fatal("Get cpu usage error")
    }
}
