package main

import (
	"os"
	"fmt"
)

var (
	Token string
)

func init() {
	Token = os.Getenv("SHIBESBOT_TOKEN")
	fmt.Print("Received token: ", Token)
}

func main() {
	initDiscord(Token)
}
