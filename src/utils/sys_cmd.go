package utils

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
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
func Cpu() (float64, error) {
    pid := os.Getpid()
    cmdstr := fmt.Sprintf("ps -p %v -o %%cpu", pid)
    output, err := Syscmd(cmdstr)
    if err != nil {
        return 0, err
    }
    percentage, err := strconv.ParseFloat(strings.Trim(strings.Split(output, "\n")[1], " "), 32)
    if err != nil {
        return 0, err
    }
    return percentage, nil
}
