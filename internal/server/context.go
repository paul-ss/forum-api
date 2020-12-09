package server

import (
	"github.com/gin-gonic/gin"
	"github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/database"
	"io"
	"os"
	"strings"
)

type Context struct {
	db *database.PgxPool
	logger *config.Logger
	ginLogFile *os.File
}

func NewContext() *Context {
	config.Conf = config.NewConfig()

	logger := config.Logger{}
	logger.Init()
	config.Lg("server", "NewContext").Info("Init logger")


	dbConn := database.NewDB()
	if err := dbConn.Open(); err != nil {
		config.Lg("server", "NewContext").Fatal("Db open: ", err.Error())
	}

	config.Lg("server", "NewContext").Info("Connected to PgxPool")

	return &Context{
		db: dbConn,
		logger: &logger,
		ginLogFile: setupGinLogger(),
	}
}

func setupGinLogger() *os.File {
	switch strings.ToLower(config.Conf.Logger.GinLevel) {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	if !config.Conf.Logger.StdoutLog {
		file, err := os.OpenFile(config.Conf.Logger.GinFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			config.Lg("server", "setupGinLogger").Fatal("Failed to log to file, using default stderr")
			return nil
		}

		gin.DefaultWriter = io.MultiWriter(file)
		return file
	} else {
		return nil
	}
}

func (c *Context) Cleanup() {
	c.db.Close()
	c.logger.Cleanup()
	c.ginLogFile.Close()
}