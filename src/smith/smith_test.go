package smith

import (
    "router"
    "testing"
)

func TestIsSuspicious(t *testing.T) {
    clue := router.Message{
        FileName: "/var/log/messages",
        Size:     403178280,
    }
    if !IsSuspicious(clue) {
        t.Fatal("/var/log/messages should be suspicious!")
    }
}
