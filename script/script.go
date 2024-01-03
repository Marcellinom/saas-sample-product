package main

import (
	"bitbucket.org/dptsi/its-go/script"
	"its.ac.id/base-go/providers"
)

func main() {
	s := script.NewScriptService()
	script.LoadFrameworkScripts(s)
	providers.LoadCustomScripts(s)

	if err := s.Run(); err != nil {
		panic(err)
	}
}
