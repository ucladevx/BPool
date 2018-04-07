package bpool

import (
	"fmt"
	"os"

	"github.com/ucladevx/BPool/adapters/http"
	"github.com/ucladevx/BPool/services"
	"github.com/ucladevx/BPool/stores/postgres"
	"github.com/ucladevx/BPool/utils/auth"

	"github.com/codyleyhan/config-loader"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

const logo = `
________  ________  ________  ________  ___          
|\   __  \|\   __  \|\   __  \|\   __  \|\  \         
\ \  \|\ /\ \  \|\  \ \  \|\  \ \  \|\  \ \  \        
 \ \   __  \ \   ____\ \  \\\  \ \  \\\  \ \  \       
  \ \  \|\  \ \  \___|\ \  \\\  \ \  \\\  \ \  \____  
   \ \_______\ \__\    \ \_______\ \_______\ \_______\
    \|_______|\|__|     \|_______|\|_______|\|_______|`

func Start() {
	env := os.Getenv("ENV")

	var conf config.LoadedData

	if env != "PROD" {
		conf = config.LoadConfig("./config", "config")
	} else {
		conf = config.LoadConfig("/config", "config")
	}

	// create logger
	loggerUnsugared, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Logger could not be created")
		return
	}
	defer loggerUnsugared.Sync()

	logger := NewBPoolLogger(loggerUnsugared.Sugar())

	// create google authorizer
	authorizer := auth.NewGoogleAuthorizer(logger)

	// create tokenizer
	tokenizer := auth.NewTokenizer(
		conf.Get("jwt.secret"),
		conf.Get("jwt.issuer"),
		int(conf.GetInt("jwt.num_days_valid")),
		logger,
	)

	// connect to db
	db := postgres.NewConnection(
		conf.Get("db.user"),
		conf.Get("db.password"),
		conf.Get("db.name"),
		conf.Get("db.port"),
		conf.Get("db.host"),
	)

	userStore := postgres.NewUserStore(db)
	postgres.CreateTables(userStore)
	userService := services.NewUserService(userStore, authorizer, tokenizer, logger)
	userController := http.NewUserController(userService, int(conf.GetInt("jwt.num_days_valid")), conf.Get("jwt.cookie"), logger)
	pagesController := http.NewPagesController(logger)

	app := echo.New()
	app.HTTPErrorHandler = handleError(logger)

	app.HideBanner = true
	app.Debug = true

	app.Pre(middleware.RequestID())
	app.Use(middleware.Logger())
	app.Use(middleware.Gzip())
	app.Use(middleware.Secure())
	app.Use(middleware.Recover())
	app.Use(auth.NewJWTmiddleware(tokenizer, conf.Get("jwt.cookie"), logger))

	fmt.Println(logo)

	pagesController.MountRoutes(app.Group(""))

	auth := app.Group("/api/v1")

	userController.MountRoutes(auth)

	port := ":" + conf.Get("port")
	logger.Info("PORT", "port", port)
	app.Logger.Fatal(app.Start(port))
}

func handleError(l *Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := 500
		message := err.Error()

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			switch v := he.Message.(type) {
			case string:
				message = v
			}
		}
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		e := c.JSON(code, echo.Map{"error": message, "request_id": requestID})

		if e != nil {
			l.Error("Handling Error, something really went wrong ", "error", err.Error())
		}

		l.Error("Handling Error ", "error", err.Error(), "status code", code)
	}
}
