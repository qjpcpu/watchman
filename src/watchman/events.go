package watchman

import (
    "strings"
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

func String(events uint32) string {
    events = events & IN_ALL_EVENTS
    list := ""
    if events == 0x0 {
        return list
    }
    if events&IN_ACCESS != 0x0 {
        list += "IN_ACCESS "
    }
    if events&IN_ATTRIB != 0x0 {
        list += "IN_ATTRIB "
    }
    if events&IN_CLOSE != 0x0 {
        if events&IN_CLOSE_NOWRITE != 0x0 {
            list += "IN_CLOSE_NOWRITE "
        } else if events&IN_CLOSE_WRITE != 0x0 {
            list += "IN_CLOSE_WRITE "
        } else {
            list += "IN_CLOSE "
        }
    }
    if events&IN_CREATE != 0x0 {
        list += "IN_CREATE "
    }
    if events&IN_DELETE != 0x0 {
        list += "IN_DELETE "
    }
    if events&IN_DELETE_SELF != 0x0 {
        list += "IN_DELETE_SELF "
    }
    if events&IN_MODIFY != 0x0 {
        list += "IN_MODIFY "
    }
    if events&IN_OPEN != 0x0 {
        list += "IN_OPEN "
    }
    if events&IN_MOVE != 0x0 {
        if events&IN_MOVED_FROM != 0x0 {
            list += "IN_MOVED_FROM "
        } else if events&IN_MOVED_TO != 0x0 {
            list += "IN_MOVED_TO "
        } else if events&IN_MOVE_SELF != 0x0 {
            list += "IN_MOVE_SELF "
        } else {
            list += "IN_MOVE "
        }
    }
    if events&IN_ISDIR != 0x0 {
        list += "IN_ISDIR "
    }
    if events&IN_IGNORED != 0x0 {
        list += "IN_IGNORED "
    }
    if events&IN_Q_OVERFLOW != 0x0 {
        list += "IN_Q_OVERFLOW "
    }
    if events&IN_UNMOUNT != 0x0 {
        list += "IN_UNMOUNT "
    }
    return strings.TrimRight(list, " ")
}
