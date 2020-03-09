package main

import (
	"goseed/utils/echoutil"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"xorm.io/core"
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

	// Configuration
	conf := viper.New()
	conf.SetConfigFile("./appsettings.json")
	if err := conf.ReadInConfig(); err != nil {
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
	tMapper := core.NewPrefixMapper(core.SnakeMapper{}, "t_")
	engine.SetTableMapper(tMapper)
	engine.SetColumnMapper(core.GonicMapper{})

	// Localization
	i18nBundle := i18n.NewBundle(language.English)
	i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	i18nBundle.MustLoadMessageFile("./utils/i18nutil/active.en.toml") // English
	i18nBundle.MustLoadMessageFile("./utils/i18nutil/active.id.toml") // Indonesian
	// i18nBundle.MustLoadMessageFile("./utils/i18nutil/active.zh.toml") // Chinese

	// Http Endpoint
	e := echo.New()
	e.Validator = &echoutil.CustomValidator{Validator: validator.New()}
	// get casbin enforcer
	enforcer := NewCasbinEnforcer(engine)
	// init all endpoints, middleware
	r := Router{
		echo:       e,
		db:         engine,
		enforcer:   enforcer,
		i18nBundle: i18nBundle,
	}
	r.Init()

	// Start server
	e.Logger.Fatal(e.Start(":8888"))
}
