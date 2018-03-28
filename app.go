package bpool

import (
	"fmt"
	"os"

	"github.com/ucladevx/BPool/adapters/http"
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

	loggerUnsugared, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Logger could not be created")
		return
	}
	defer loggerUnsugared.Sync()

	logger := loggerUnsugared.Sugar()

	authorizer := auth.NewGoogleAuthorizer(
		conf.Get("google.id"),
		conf.Get("google.secret"),
		conf.Get("google.redirect_url"),
		conf.Get("secret"),
		logger,
	)
	authController := http.NewAuthController(authorizer, logger)

	app := echo.New()
	app.HTTPErrorHandler = handleError(logger)

	app.HideBanner = true
	app.Debug = true

	app.Pre(middleware.RequestID())
	app.Use(middleware.Logger())
	app.Use(middleware.Gzip())
	app.Use(middleware.Secure())
	app.Use(middleware.Recover())

	fmt.Println(logo)

	app.GET("/", func(c echo.Context) error {
		logger.Infow("INDEX ROUTE", "request id", "test")
		return c.HTML(200, "<html><title>Golang Google</title> <body> <a href='/api/auth/login'><button>Login with Google!</button> </a> </body></html>")
	})

	auth := app.Group("/api/auth")

	authController.MountRoutes(auth)

	port := ":" + conf.Get("port")
	logger.Infow("PORT", "port", port)
	app.Logger.Fatal(app.Start(port))
}

func handleError(l *zap.SugaredLogger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := 500
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		e := c.JSON(code, echo.Map{"error": err.Error()})

		if e != nil {
			l.Error("Handling Error, something really went wrong", "error", err.Error())
		}

		l.Error("Handling Error", "error", err.Error(), "status code", code)
	}
}
