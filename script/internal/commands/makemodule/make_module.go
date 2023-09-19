package makemodule

import (
	"fmt"
	"os"

	"github.com/mikestefanello/hooks"
	"its.ac.id/base-go/script/internal/app"
)

func init() {
	app.HookBoot.Listen(func(event hooks.Event[*app.Script]) {
		event.Msg.AddCommand(app.Command{
			Name:        "make:module",
			Description: "Create new module",
			Handler:     makeModule,
		})
	})
}

func makeModule(args []string) {
	if len(args) == 0 {
		fmt.Println("No module name provided")
		return
	}
	name := args[0]
	path := fmt.Sprintf("modules/%s", name)

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Module %s already exist\n", name)
		return
	}

	os.MkdirAll(path, os.ModePerm)
	createSkeleton(name, path)
}

func getBasePkgPath() (string, error) {
	goModFile, err := os.Open("go.mod")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var basePkgPath string
	fmt.Fscanf(goModFile, "module %s", &basePkgPath)
	goModFile.Close()

	return basePkgPath, nil
}

func createSkeleton(name string, path string) {
	basePkgPath, err := getBasePkgPath()
	if err != nil {
		fmt.Println(err)
		return
	}

	moduleInitFile, err := os.Create(fmt.Sprintf("%s/%s.go", path, name))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(
		moduleInitFile,
		`package %s
		
import _ "%s/modules/%s/internal/services/routes"

func init() {

}
		`,
		name,
		basePkgPath,
		name,
	)
	moduleInitFile.Close()

	err = os.MkdirAll(fmt.Sprintf("%s/internal/services/routes", path), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	moduleRoutesFile, err := os.Create(fmt.Sprintf("%s/internal/services/routes/routes.go", path))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(
		moduleRoutesFile,
		`package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
	"%s/services/web"
)

func registerRoutes(r *gin.Engine) {
	g := r.Group("/%s")
	i := do.DefaultInjector

	// Register routes below

}

func init() {
	web.HookBuildRouter.Listen(func(event hooks.Event[*gin.Engine]) {
		registerRoutes(event.Msg)
	})
}
		`,
		basePkgPath,
		name,
	)
	moduleRoutesFile.Close()

	var moduleFolders = []string{
		"internal/app",
		"internal/app/controllers",
		"internal/app/commands",
		"internal/app/queries",
		"internal/app/services",

		"internal/domain",
		"internal/domain/entities",
		"internal/domain/events",
		"internal/domain/repositories",
		"internal/domain/services",
		"internal/domain/valueobjects",

		"internal/infrastructures",
	}

	for _, folder := range moduleFolders {
		os.MkdirAll(fmt.Sprintf("%s/%s", path, folder), os.ModePerm)
		gitKeepFile, err := os.Create(fmt.Sprintf("%s/%s/.gitkeep", path, folder))
		if err != nil {
			fmt.Println(err)
			return
		}
		gitKeepFile.Close()
	}
}
