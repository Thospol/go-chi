package main

import (
	"flag"
	"saaa-api/internal/core/config"
	"saaa-api/internal/core/jwt"
	"saaa-api/internal/core/sql"
	"saaa-api/internal/handlers/routes"
	"saaa-api/internal/handlers/servers"
	"strconv"

	"saaa-api/docs"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	environment := flag.String("environment", "local", "set working environment")
	configs := flag.String("config", "configs", "set configs path, default as: 'configs'")

	flag.Parse()

	// Init configuration
	if err := config.InitConfig(*configs, *environment); err != nil {
		panic(err)
	}
	//=======================================================

	// programatically set swagger info
	docs.SwaggerInfo.Title = config.CF.SwaggerInfo.Title
	docs.SwaggerInfo.Description = config.CF.SwaggerInfo.Description
	docs.SwaggerInfo.Version = config.CF.SwaggerInfo.Version
	docs.SwaggerInfo.Host = config.CF.SwaggerInfo.Host
	//=======================================================

	// set logrus
	if config.CF.App.Release {
		logrus.SetFormatter(stackdriver.NewFormatter(
			stackdriver.WithService("api"),
			stackdriver.WithVersion("v1.0.0")))
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	logrus.Infof("Initial 'Configuration'. %+v", config.CF)
	//=======================================================

	// Init return result
	if err := config.InitReturnResult("configs"); err != nil {
		panic(err)
	}
	//=======================================================

	// Creates a new google cloud storage client
	// storage.NewClient()
	// =======================================================

	// Get SecretKey JWT
	jwt.LoadKey()
	// =======================================================

	// Init connection postgresql
	mysqlConfig, err := sql.InitConnectionMysql(config.CF)
	if err != nil {
		panic(err)
	}
	mysqlConfig.ExportDatabase()
	//========================================================

	// NewRouter && NewServer
	r := routes.NewRouter()
	srv := servers.NewServer(strconv.Itoa(config.CF.App.Port), r)
	srv.ListenAndServeWithGracefulShutdown()
	//=======================================================
}
