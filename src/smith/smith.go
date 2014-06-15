package smith

import (
    "container/list"
    "math"
    . "mlog"
    "router"
    "syscall"
    "time"
    "watchman"
)

const (
    SizeLimitPercentage = 0.001
    MonthLimitAgo       = 3
)

func IsSuspicious(clue router.Message) bool {
    file := clue.FileName
    fs := syscall.Statfs_t{}
    syscall.Statfs(file, &fs)
    if percentage := float64(clue.Size) / (float64(fs.Bsize) * float64(fs.Blocks)); !math.IsInf(percentage, 1) && percentage > SizeLimitPercentage {
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
        for {
            if ele := queue.Back(); ele != nil {
                value := queue.Remove(ele)
                msg := value.(router.Message)
                if IsSuspicious(msg) {
                    Log.Debugf("I will kill %s(%s)", msg.FileName, watchman.HumanReadable(msg.Event))
                } else {
                    Log.Debugf("You're good %s(%s), let you go.", msg.FileName, watchman.HumanReadable(msg.Event))
                }
            } else {
                break
            }
        }
    }
}
