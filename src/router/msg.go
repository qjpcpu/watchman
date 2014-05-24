package router

import (
    "encoding/json"
    "fmt"
)

type Message struct {
    Event                                        uint32
    FileName, AccessTime, ModifyTime, ChangeTime string
    Inode, Size                                  int
}

func ParseMessage(content string) (Message, error) {
    var m Message
    err := json.Unmarshal([]byte(content), &m)
    return m, err
}
func (m Message) String() string {
    format := `
    {
        "Event":%v,
        "FileName":"%v",
        "AccessTime":"%v",
        "ModifyTime":"%v",
        "ChangeTime":"%v",
        "Inode":%v,
        "Size":%v
    }
    `
    return fmt.Sprintf(format, m.Event, m.FileName, m.AccessTime, m.ModifyTime, m.ChangeTime, m.Inode, m.Size)
}
