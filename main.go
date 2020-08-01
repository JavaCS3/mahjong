package main

import (
	"fmt"
	"log"
	"os/exec"

	"mahjong.com/pkg/core"
	"mahjong.com/pkg/utils"
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func main() {
	c := exec.Command("node", "mytile/index.js")
	stdout, err := c.StdoutPipe()
	checkError(err)
	go utils.ScanLines(
		stdout,
		func(s string) error {
			c, err := core.CmdFromStr(s)
			if err == nil {
				fmt.Printf("[CMD] %v\n", c)
			} else {
				fmt.Println(s)
			}
			return nil
		},
	)
	c.Run()
}
