package router

import (
    "testing"
    "time"
)

func TestStart(t *testing.T) {
    Start(nil)
    time.Sleep(time.Second * 2)
}
