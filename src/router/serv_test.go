package router

import (
    "testing"
)

func TestStart(t *testing.T) {
    Start(DefaultBuilder())
    c, err := NewRouterCli(SYS_ID, defaultSocketFunc)
    if err != nil {
        t.Fatal(err)
    }
    c1, err := NewRouterCli("/", defaultSocketFunc)
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
    c1.Close()
    c.Write("end")
    err = c1.Write("x")
    if err.Error() != "use of closed network connection" {
        t.Fatal("the connection should be closed!")
    }
    c2, err := NewRouterCli("CLIENT-2", defaultSocketFunc)
    if err != nil {
        t.Fatal(err)
    }
    c2.Write("I'm client-2")
    if str, err := c.Read(); err != nil {
        t.Fatal("should read a message")
    } else if str != "I'm client-2" {
        t.Fatalf("should got %v", "I'm client-2")
    }

}
