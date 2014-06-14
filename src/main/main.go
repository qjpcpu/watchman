package main

import (
    "alfred"
    "container/list"
    . "mlog"
    "os"
    "os/signal"
    "smith"
    "syscall"
    "watchman"
)

func bigWatch() {
    man, err := watchman.NewWatchman()
    if err != nil {
        Log.Fatal(err)
    }
    wlist := getWatchlist()
    for _, f := range wlist {
        if err = man.WatchPath(f, watchman.IN_ALL_EVENTS); err != nil {
            Log.Errorf("%s: %v", f, err)
        }
    }
    go func() {
        events := list.New()
        go smith.ScanAbnormal(events)
        for {
            if m, err := man.PullEvent(); err == nil {
                events.PushFront(m)
            }
        }
    }()
}
func configLogger() {
    SetLevel("DEBUG")
}
func main() {
    configLogger()
    alfred.Boot()
    defer alfred.Shutdown()
    bigWatch()
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
    <-sigc
}
