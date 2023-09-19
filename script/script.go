package main

import (
	"fmt"
	"os"

	"its.ac.id/base-go/script/internal/app"

	// Commands
	_ "its.ac.id/base-go/script/internal/commands/makemodule"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No arguments provided")
		return
	}

	s := app.NewScript()
	app.HookBoot.Dispatch(s)

	command, exist := s.GetCommand(args[0])
	if !exist {
		fmt.Printf("Command %s not found\n", args[0])
		return
	}
	command.Handler(args[1:])
}
