package main

import (
    _ "alfred"
    . "mlog"
    "os"
    "os/signal"
    "syscall"
    "utils"
    "watchman"
)

func bigWatch() {
    man, err := watchman.NewWatchman()
    if err != nil {
        Log.Fatal(err)
    }
    list := utils.Walk("/home", 5)
    list = list[0:len(list)]
    for _, f := range list {

        if err = man.WatchPath(f, watchman.IN_ALL_EVENTS); err != nil {
            Log.Info("ERROR", err)
        }
    }
    go func() {
        for {
            if m, err := man.PullEvent(); err == nil {
                Log.Info(m.FileName, watchman.String(m.Event))
            }
        }
    }()
}
func configLogger() {
    SetLevel("INFO")
}
func main() {
    configLogger()
    Log.Info("START")
    bigWatch()
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
    <-sigc
    Log.Info("Shutting down.")
}
