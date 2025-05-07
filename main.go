package main

import (
	"fmt"
	_ "github.com/EnderCHX/go-qq-aibot/config"
	_ "github.com/EnderCHX/go-qq-aibot/qqbot"
	"os"
)

func main() {
	items, err := os.ReadFile("items.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(items))
}
