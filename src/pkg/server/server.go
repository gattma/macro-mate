package server

import (
	"fmt"
	"net/http"

	"github.com/gattma/macro-mate/src/cmd/utils"
	"github.com/gattma/macro-mate/src/pkg/client/mongodb"
	"github.com/gattma/macro-mate/src/pkg/handler"
	"github.com/gattma/macro-mate/src/pkg/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Initialize(config utils.Configuration) {

	// Create a new instance of the logger. You can have any number of instances.
	var log = logrus.New()

	log.WithFields(logrus.Fields{
		"mongo_url":   config.Database.Url,
		"server_port": config.Server.Port,
		"db_name":     config.Database.DbName,
		"collection":  config.Database.Collection,
		"timeout":     config.App.Timeout,
	}).Info("\nConfiguration informations\n")

	logrus.Infof("Application Name %s is starting....", config.App.Name)

	client, err := mongodb.ConnectMongoDb(config.Database.Url)

	if err != nil {
		logrus.Fatal(err)
	}

	repo := repository.NewFoodRepository(&config, client)
	macroHandler := handler.NewMacroMateHandler(client, repo, config)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.Use(cors.Default())

	router.LoadHTMLGlob("web/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Macro Mate",
		})
	})

	// Route for about page
	router.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", gin.H{
			"title": "Neues Produkt",
		})
	})

	api := router.Group("api/v1")
	{
		api.GET("/health", macroHandler.Healthcheck)
		api.GET("/macromate/search", macroHandler.Search)
		api.GET("/macromate/searchkeyword", macroHandler.SearchKeywords)
		api.POST("/macromate/add", macroHandler.Add)
	}

	// PORT environment variable was defined.
	formattedUrl := fmt.Sprintf(":%s", config.Server.Port)

	router.Run(formattedUrl)
}
