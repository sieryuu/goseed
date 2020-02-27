package main

import (
	"go.uber.org/zap"
	"xorm.io/xorm"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	logger *zap.Logger
)

func main() {
	// init zap
	logger, _ = zap.NewProduction()
	defer logger.Sync()

	conf := viper.New()
	conf.SetConfigFile("./appsettings.json")
	err := conf.ReadInConfig()
	if err != nil {
		logger.Error("failed to init viper", zap.Error(err))
		panic(err)
	}

	// Setup Database
	engine, err := xorm.NewEngine("postgres", conf.GetString("connectionstrings.default"))
	defer engine.Close()

	if err != nil {
		logger.Error("failed to init database", zap.Error(err))
		panic(err)
	}

	// Http Endpoint
	e := echo.New()
	// get casbin enforcer
	enforcer := NewCasbinEnforcer(engine)
	// init all endpoints, middleware
	ConfigureEcho(e, engine, enforcer)

	// Start server
	e.Logger.Fatal(e.Start(":8888"))
}
