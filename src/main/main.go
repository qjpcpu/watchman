package main

import (
    _ "alfred"
    "log"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    log.Println("START")
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
    <-sigc
    log.Println("Shutting down.")
}
