package main

import (
    "utils"
)

func getWatchlist() (list []string) {
    list = utils.Walk("/var/log", 100)
    list = append(list, utils.Walk("/home", 10)...)
    return
}
