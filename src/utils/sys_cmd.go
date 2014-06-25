package utils

import (
    "os/exec"
    "strings"
)

func Syscmd(cmdstr string) (string, error) {
    cmd := strings.Fields(cmdstr)
    out, err := exec.Command(cmd[0], cmd[1:]...).Output()
    if err != nil {
        return "", err
    } else {
        return string(out), nil
    }
}
