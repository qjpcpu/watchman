package main

import (
    _ "alfred"
    "log"
    "os"
    "os/signal"
    "syscall"
    "watchman"
)

func createWatcher() {
    man, err := watchman.NewWatchman()
    if err != nil {
        log.Fatal(err)
    }
    if err = man.WatchPath("/tmp", watchman.IN_ALL_EVENTS); err != nil {
        log.Println(err)
    }
    go func() {
        for {
            if m, err := man.PullEvent(); err == nil {
                log.Println(m)
            } else if err.Error() == "SYSTEM" {
                log.Println(m.FileName)
            }
        }
        man.Release()
    }()
    //if err = man.ForgetPath("/home/work/repository/watchman/src/watchman"); err != nil {
    //    log.Println(err)
    //}
}
func main() {
    log.Println("START")
    createWatcher()
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
    <-sigc
    log.Println("Shutting down.")
}
