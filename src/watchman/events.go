package watchman

import (
    "syscall"
)

const (

    // Options for AddWatch
    IN_DONT_FOLLOW uint32 = syscall.IN_DONT_FOLLOW
    IN_ONESHOT     uint32 = syscall.IN_ONESHOT
    IN_ONLYDIR     uint32 = syscall.IN_ONLYDIR

    // Events
    IN_ACCESS        uint32 = syscall.IN_ACCESS
    IN_ALL_EVENTS    uint32 = syscall.IN_ALL_EVENTS
    IN_ATTRIB        uint32 = syscall.IN_ATTRIB
    IN_CLOSE         uint32 = syscall.IN_CLOSE
    IN_CLOSE_NOWRITE uint32 = syscall.IN_CLOSE_NOWRITE
    IN_CLOSE_WRITE   uint32 = syscall.IN_CLOSE_WRITE
    IN_CREATE        uint32 = syscall.IN_CREATE
    IN_DELETE        uint32 = syscall.IN_DELETE
    IN_DELETE_SELF   uint32 = syscall.IN_DELETE_SELF
    IN_MODIFY        uint32 = syscall.IN_MODIFY
    IN_MOVE          uint32 = syscall.IN_MOVE
    IN_MOVED_FROM    uint32 = syscall.IN_MOVED_FROM
    IN_MOVED_TO      uint32 = syscall.IN_MOVED_TO
    IN_MOVE_SELF     uint32 = syscall.IN_MOVE_SELF
    IN_OPEN          uint32 = syscall.IN_OPEN

    // Special events
    IN_ISDIR      uint32 = syscall.IN_ISDIR
    IN_IGNORED    uint32 = syscall.IN_IGNORED
    IN_Q_OVERFLOW uint32 = syscall.IN_Q_OVERFLOW
    IN_UNMOUNT    uint32 = syscall.IN_UNMOUNT
)
