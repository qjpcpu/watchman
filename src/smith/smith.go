package smith

import (
    "container/list"
    . "mlog"
    "router"
    "syscall"
    "time"
)

const (
    SizeLimitPercentage = 0.001
    MonthLimitAgo       = 3
)

func IsSuspicious(clue router.Message) bool {
    file := clue.FileName
    fs := syscall.Statfs_t{}
    syscall.Statfs(file, &fs)
    if percentage := float32(clue.Size) / (float32(fs.Bsize) * float32(fs.Blocks)); percentage > SizeLimitPercentage {
        return true
    }
    return false
}

func FindSuspicious(dir string) (list []string) {

    return
}

func ScanAbnormal(queue *list.List) {
    c := time.Tick(1 * time.Second)
    for _ = range c {
        if ele := queue.Back(); ele != nil {
            queue.Remove(ele)
            msg := ele.Value.(router.Message)
            if IsSuspicious(msg) {
                Log.Debug("I will kill", msg.FileName)
            } else {
                Log.Debugf("You're good %s, let you go.", msg.FileName)
            }
        }
    }
}
