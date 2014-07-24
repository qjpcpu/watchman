package smith

import (
    "alfred"
    "container/list"
    "github.com/qjpcpu/logger"
    "time"
    "utils"
)

func ScanAbnormal(queue *list.List) {
    c := time.Tick(time.Second)
    action := "info"
    if mcfg, err := utils.GetMainConfig(); err == nil {
        action = mcfg.Action
    }
    for _ = range c {
        erase_list := []alfred.Message{}
        for {
            if ele := queue.Back(); ele != nil {
                value := queue.Remove(ele)
                msg := value.(alfred.Message)
                if fromWhiteList(msg) {
                    logger.LoggerOf("watchman-logger").Infof("%v is on  white list,pass.", msg.FileName)
                } else if fromBigFile(msg) {
                    if canEraseInstant(msg.FileName) {
                        erase_list = append(erase_list, msg)
                    } else {
                        logger.LoggerOf("watchman-logger").Warningf("Big file found and I dare not del.(%v:%v)", msg.FileName, msg.Size)
                    }
                } else if can_del, ok := fromBigDirectory(msg); ok {
                    var candel []alfred.Message
                    for _, m := range can_del {
                        candel = append(candel, alfred.Message{FileName: m})
                    }
                    if yes, _ := canErase(candel...); len(yes) > 0 {
                        erase_list = append(erase_list, yes...)
                    }
                } else {
                    logger.LoggerOf("watchman-logger").Debugf("You're good %s(%s), let you go.", msg.FileName, alfred.HumanReadable(msg.Event))
                }
            } else {
                break
            }
        }
        if len(erase_list) > 0 {
            logger.LoggerOf("watchman-logger").Infof("[%s] Remove %v", action, erase_list)
            switch action {
            case "info":
                printState(erase_list...)
            case "remove":
                erase(erase_list...)
            }
        }
    }
}
