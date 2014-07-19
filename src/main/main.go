package main

import (
    "alfred"
    "bitbucket.org/kardianos/osext"
    "container/list"
    "fmt"
    "github.com/qjpcpu/go-logging"
    "os"
    "os/signal"
    "path/filepath"
    "smith"
    "strings"
    "syscall"
    "utils"
    "watchman"
)

const (
    ICARE_EVENTS = watchman.IN_ISDIR | watchman.IN_CLOSE_WRITE | watchman.IN_CREATE | watchman.IN_DELETE | watchman.IN_DELETE_SELF | watchman.IN_MOVE | watchman.IN_MODIFY
)

func bigWatch() {
    man, err := watchman.NewWatchman()
    if err != nil {
        logging.Fatal(err)
    }
    wlist := utils.GetWatchlist()
    for _, f := range wlist {
        if err = man.WatchPath(f, ICARE_EVENTS); err != nil {
            logging.Errorf("%s: %v", f, err)
        }
    }
    go func() {
        events := list.New()
        go smith.ScanAbnormal(events)
        for {
            if m, err := man.PullEvent(); err == nil {
                if m.Event&watchman.IN_ISDIR != 0 && m.Event&watchman.IN_CREATE != 0 {
                    man.WatchPath(m.FileName, ICARE_EVENTS)
                } else if m.Event&watchman.IN_ISDIR != 0 && m.Event&watchman.IN_DELETE != 0 {
                    man.ForgetPath(m.FileName)
                } else {
                    events.PushFront(m)
                }
            }
        }
    }()
}
func configLogger() {
    cfg, err := utils.GetMainConfig()
    if err == nil {
        level := logging.INFO
        if cfg.LogLevel != "" {
            switch strings.ToUpper(cfg.LogLevel) {
            case "DEBUG":
                level = logging.DEBUG
            case "INFO":
                level = logging.INFO
            case "CRITICAL":
                level = logging.CRITICAL
            case "ERROR":
                level = logging.ERROR
            case "WARNING":
                level = logging.WARNING
            case "NOTICE":
                level = logging.NOTICE
            }
        }
        if cfg.LogFile != "" {
            logging.InitSimpleFileLogger(cfg.LogFile, level)
        } else {
            logging.InitLogger(level)
        }
        return
    }
    logging.InitLogger(logging.DEBUG)
}
func writePidfile() {
    if filename, err := osext.Executable(); err == nil {
        pid := fmt.Sprintf("%s/%s.pid", filepath.Dir(filepath.Dir(filename)), filepath.Base(filename))
        if fi, err := os.Create(pid); err == nil {
            fi.Write([]byte(fmt.Sprintf("%v", os.Getpid())))
            defer fi.Close()
        }
    }
}
func main() {
    writePidfile()
    configLogger()
    alfred.Boot()
    defer alfred.Shutdown()
    bigWatch()
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
    <-sigc
}
