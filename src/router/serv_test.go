package router

import (
    "testing"
)

func TestStart(t *testing.T) {
    Start(defaultPolicy)
    c, err := NewRouterCli(SYS_ID)
    if err != nil {
        t.Fatal(err)
    }
    c1, err := NewRouterCli("/")
    if err != nil {
        t.Fatal(err)
    }
    c.Write("hello all")
    c1.Write("hello alfred")
    if str, err := c1.Read(); err != nil {
        t.Fatal("should read a message")
    } else if str != "hello all" {
        t.Fatalf("should got %v", "hello all")
    }
    if str, err := c.Read(); err != nil {
        t.Fatal("should read a message")
    } else if str != "hello alfred" {
        t.Fatalf("should got %v", "hello alfred")
    }
}
