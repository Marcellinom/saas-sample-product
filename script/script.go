package main

import (
	"bitbucket.org/dptsi/its-go/script"
)

func main() {
	s := script.NewScriptService()

	if err := s.Run(); err != nil {
		panic(err)
	}
}
