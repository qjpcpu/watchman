package main

import (
    _ "alfred"
    "fmt"
)

func main() {
    fmt.Println("start...")
    var word []byte
    fmt.Scanln(&word)
    fmt.Println(string(word))
}
