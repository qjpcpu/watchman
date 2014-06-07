package mlog

import (
    "fmt"
    "github.com/op/go-logging"
)

var log_id = "mlog"

var Log *mlogger

type mlogger struct {
    *logging.Logger
}

func (mlog *mlogger) Critical(args ...interface{}) {
    mlog.Logger.Critical("%s", fmt.Sprintln(args...))
}
func (mlog *mlogger) Criticalf(format string, args ...interface{}) {
    mlog.Logger.Critical("%s", fmt.Sprintf(format, args...))
}
func (mlog *mlogger) Debug(args ...interface{}) {
    mlog.Logger.Debug("%s", fmt.Sprintln(args...))
}
func (mlog *mlogger) Debugf(format string, args ...interface{}) {
    mlog.Logger.Debug("%s", fmt.Sprintf(format, args...))
}
func (mlog *mlogger) Error(args ...interface{}) {
    mlog.Logger.Error("%s", fmt.Sprintln(args...))
}
func (mlog *mlogger) Errorf(format string, args ...interface{}) {
    mlog.Logger.Error("%s", fmt.Sprintf(format, args...))
}
func (mlog *mlogger) Info(args ...interface{}) {
    mlog.Logger.Info("%s", fmt.Sprintln(args...))
}
func (mlog *mlogger) Infof(format string, args ...interface{}) {
    mlog.Logger.Info("%s", fmt.Sprintf(format, args...))
}
func (mlog *mlogger) Notice(args ...interface{}) {
    mlog.Logger.Notice("%s", fmt.Sprintln(args...))
}
func (mlog *mlogger) Noticef(format string, args ...interface{}) {
    mlog.Logger.Notice("%s", fmt.Sprintf(format, args...))
}
func (mlog *mlogger) Warning(args ...interface{}) {
    mlog.Logger.Warning("%s", fmt.Sprintln(args...))
}
func (mlog *mlogger) Warningf(format string, args ...interface{}) {
    mlog.Logger.Warning("%s", fmt.Sprintf(format, args...))
}
func init() {
    format := logging.MustStringFormatter("%{level} %{message}")
    logging.SetFormatter(format)
    innerlog := logging.MustGetLogger(log_id)
    logging.SetLevel(logging.INFO, log_id)
    Log = &mlogger{innerlog}
}

func SetLevel(level string) {
    switch level {
    case "CRITICAL":
        logging.SetLevel(logging.CRITICAL, log_id)
    case "ERROR":
        logging.SetLevel(logging.ERROR, log_id)
    case "WARNING":
        logging.SetLevel(logging.WARNING, log_id)
    case "NOTICE":
        logging.SetLevel(logging.NOTICE, log_id)
    case "INFO":
        logging.SetLevel(logging.INFO, log_id)
    case "DEBUG":
        logging.SetLevel(logging.DEBUG, log_id)
    }
}
