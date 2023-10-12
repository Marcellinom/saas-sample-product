package makemodule

import (
	"fmt"
	"os"
	"strings"

	"github.com/mikestefanello/hooks"
	"github.com/stoewer/go-strcase"
	"its.ac.id/base-go/script/internal/app"
)

func init() {
	app.HookBoot.Listen(func(event hooks.Event[*app.Script]) {
		event.Msg.AddCommand(app.Command{
			Name:        "make:module",
			Description: "Create new module",
			Usage:       "make:module <module_name>",
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

	var ans string
	for {
		fmt.Print("Do you want to use transaction script instead of aggregate pattern? (y/N): ")
		fmt.Scanln(&ans)
		if ans == "y" || ans == "N" || ans == "" {
			break
		}
		fmt.Println("Invalid answer")
	}
	tsPattern := ans == "y"

	os.MkdirAll(path, os.ModePerm)
	if err := createSkeleton(name, path, tsPattern); err != nil {
		fmt.Println(err)
		return
	}
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

func createSkeleton(name string, path string, tsPattern bool) error {
	basePkgPath, err := getBasePkgPath()
	if err != nil {
		return err
	}

	if err := createModuleInitFile(path, name, basePkgPath); err != nil {
		return err
	}

	if err := createModuleFolders(path, tsPattern); err != nil {
		return err
	}

	if err := createRoutesFile(path, basePkgPath, name); err != nil {
		return err
	}

	if err := createListenersFile(path, basePkgPath, name); err != nil {
		return err
	}

	if err := createConfigFile(path, basePkgPath, name); err != nil {
		return err
	}

	return nil
}

func createModuleFolders(path string, tsPattern bool) error {
	var moduleFolders = []string{
		"internal/app/config",
		"internal/app/commands",
		"internal/app/listeners",
		"internal/app/queries",
		"internal/app/routes",
		"internal/app/services",

		"internal/infrastructures/database",

		"internal/domain/services",
		"internal/domain/repositories",

		"internal/presentation/controllers",
	}

	if !tsPattern {
		moduleFolders = append(
			moduleFolders,
			"internal/domain/entities",
			"internal/domain/events",
			"internal/domain/valueobjects",
		)
	}

	for _, folder := range moduleFolders {
		os.MkdirAll(fmt.Sprintf("%s/%s", path, folder), os.ModePerm)
		gitKeepFile, err := os.Create(fmt.Sprintf("%s/%s/.gitkeep", path, folder))
		if err != nil {
			return err
		}
		gitKeepFile.Close()
	}
	return nil
}

func createRoutesFile(path string, basePkgPath string, name string) error {
	moduleRoutesFile, err := os.Create(fmt.Sprintf("%s/internal/app/routes/routes.go", path))
	if err != nil {
		return err
	}
	pascalcased := strcase.UpperCamelCase(name)

	dsn := "fmt.Sprintf(\"sqlserver://%s:%s@%s:%s?database=%s\", dbCfg.SQLServerUsername, dbCfg.SQLServerPassword, dbCfg.SQLServerHost, dbCfg.SQLServerPort, dbCfg.SQLServerDatabase)"
	fmt.Fprintf(
		moduleRoutesFile,
		`package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
	"%s/bootstrap/web"
	"%s/modules/%s/internal/app/config"
	"%s/pkg/app"
)

const routePrefix = "/%s"

func registerRoutes(r *gin.Engine) {
	g := r.Group(routePrefix)

	// Register routes below

}

func init() {
	web.HookBuildRouter.Listen(func(event hooks.Event[*gin.Engine]) {
		registerRoutes(event.Msg)
	})

	app.HookBoot.Listen(func(event hooks.Event[*do.Injector]) {
		i := event.Msg

		cfg := do.MustInvoke[config.%sConfig](i)
		dbCfg := cfg.Database()

		// Uncomment this line if you want to use Firestore
		// client, err := firestore.NewClient(context.Background(), dbCfg.FirestoreProjectID)

		// Uncomment this line if you want to use SQL Server
		// dsn := %s
		// db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		// if err != nil {
		// 	panic("failed to connect SQL Server database for session")
		// }

		// Register services below
	})
}
		`,
		basePkgPath,
		basePkgPath,
		strings.ReplaceAll(name, "_", "-"),
		basePkgPath,
		strings.ReplaceAll(name, "_", "-"),
		pascalcased,
		dsn,
	)
	moduleRoutesFile.Close()
	return nil
}

func createListenersFile(path string, basePkgPath string, name string) error {
	moduleListenersFile, err := os.Create(fmt.Sprintf("%s/internal/app/listeners/listeners.go", path))
	if err != nil {
		return err
	}
	fmt.Fprintf(
		moduleListenersFile,
		`package listeners

import (
	"github.com/mikestefanello/hooks"
	"%s/pkg/app/common"
	"%s/bootstrap/event"
)

func registerListeners(e *common.Event) {
	
}

func init() {
	event.HookEvent.Listen(func(event hooks.Event[*common.Event]) {
		registerListeners(event.Msg)
	})
}
		`,
		basePkgPath,
		basePkgPath,
	)
	moduleListenersFile.Close()
	return nil
}

func createConfigFile(path string, basePkgPath string, name string) error {
	moduleConfigFile, err := os.Create(fmt.Sprintf("%s/internal/app/config/config.go", path))
	if err != nil {
		return err
	}
	defer moduleConfigFile.Close()
	pascalCased := strcase.UpperCamelCase(name)
	fmt.Fprintf(
		moduleConfigFile,
		`package config

import (
	"github.com/samber/do"
)

type %sConfig interface {
	Database() DatabaseConfig
}

type %sConfigImpl struct {
	db DatabaseConfig
}

func NewConfig(i *do.Injector) (%sConfig, error) {
	db := setupDatabaseConfig()

	return %sConfigImpl{db}, nil
}

func init() {
	do.Provide[%sConfig](do.DefaultInjector, NewConfig)
}
		`,
		pascalCased,
		pascalCased,
		pascalCased,
		pascalCased,
		pascalCased,
	)

	moduleDbConfigFile, err := os.Create(fmt.Sprintf("%s/internal/app/config/database.go", path))
	if err != nil {
		return err
	}
	defer moduleDbConfigFile.Close()

	uppercased := strings.ToUpper(name)
	fmt.Fprintf(moduleDbConfigFile,
		"package config\n"+
			"\n"+
			"import (\n"+
			"	\"github.com/joeshaw/envdecode\"\n"+
			")\n"+
			"\n"+
			"type DatabaseConfig struct {\n"+
			"	// Firestore session adapter\n"+
			"	FirestoreProjectID  string `env:\"%s_FIRESTORE_PROJECT_ID\"`\n"+
			"	FirestoreCollection string `env:\"%s_FIRESTORE_COLLECTION\"`\n"+
			"\n"+
			"	// SQLite session adapter (GORM)\n"+
			"	SQLiteDB string `env:\"%s_SQLITE_DB\"`\n"+
			"\n"+
			"	// SQL Server session adapter (GORM)\n"+
			"	SQLServerHost     string `env:\"%s_SQLSERVER_HOST\"`\n"+
			"	SQLServerPort     string `env:\"%s_SQLSERVER_PORT\"`\n"+
			"	SQLServerDatabase string `env:\"%s_SQLSERVER_DATABASE\"`\n"+
			"	SQLServerUsername string `env:\"%s_SQLSERVER_USERNAME\"`\n"+
			"	SQLServerPassword string `env:\"%s_SQLSERVER_PASSWORD\"`\n"+
			"}\n"+
			"\n"+
			"func (c %sConfigImpl) Database() DatabaseConfig {\n"+
			"	return c.db\n"+
			"}\n"+
			"\n"+
			"func setupDatabaseConfig() DatabaseConfig {\n"+
			"	var db DatabaseConfig\n"+
			"	err := envdecode.StrictDecode(&db)\n"+
			"	if err != nil {\n"+
			"		panic(err)\n"+
			"	}\n"+
			"\n"+
			"	return db\n"+
			"}\n",
		uppercased,
		uppercased,
		uppercased,
		uppercased,
		uppercased,
		uppercased,
		uppercased,
		uppercased,
		pascalCased,
	)

	return nil
}

func createModuleInitFile(path string, name string, basePkgPath string) error {
	moduleInitFile, err := os.Create(fmt.Sprintf("%s/%s.go", path, name))
	if err != nil {
		return err
	}
	fmt.Fprintf(
		moduleInitFile,
		`package %s
		

import (
	_ "%s/modules/%s/internal/app/config"
	_ "%s/modules/%s/internal/app/listeners"
	_ "%s/modules/%s/internal/app/routes"
)

func init() {

}
		`,
		name,
		basePkgPath,
		name,
		basePkgPath,
		name,
		basePkgPath,
		name,
	)
	moduleInitFile.Close()
	return nil
}
