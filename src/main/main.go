package main

import (
    _ "alfred"
    "log"
    "os"
    "os/signal"
    "syscall"
    "utils"
    "watchman"
)

func bigWatch() {
    man, err := watchman.NewWatchman()
    if err != nil {
        log.Fatal(err)
    }
    list := utils.Walk("/home", 100)
    list = list[0:len(list)]
    for _, f := range list {

        if err = man.WatchPath(f, watchman.IN_ALL_EVENTS); err != nil {
            log.Println("ERROR", err)
        }
    }
    go func() {
        for {
            if m, err := man.PullEvent(); err == nil {
                log.Println(m.FileName, watchman.String(m.Event))
            }
        }
    }()
}
func main() {
    log.Println("START")
    bigWatch()
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
    <-sigc
    log.Println("Shutting down.")
}
