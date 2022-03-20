package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"log"
	"whatsapp-client/api"
)

var host = flag.String("host", "127.0.0.1", "")
var port = flag.String("port", "9000", "")

func main() {
	flag.Parse()

	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Fatal error config file: %w \n", err)
	}

	g := gin.New()
	g.Use(gin.Logger())
	g.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.JSON(200, fmt.Sprint(err))
	}))
	g.Use(ResponseMiddleware())
	g.Use(cors.Default())

	group := g.Group("/api")
	group.POST("/upload", api.UploadAdd)
	group.GET("/upload", api.UploadGet)
	group.GET("/devices", api.DeviceQuery)
	group.DELETE("/device", api.DeviceDelete)
	group.GET("/connect", api.ClientLogin)
	group.GET("/disconnect", api.ClientLogout)
	group.GET("/info", api.ClientInfo)
	group.GET("/groups", api.GroupQuery)
	group.GET("/group", api.GroupGet)
	group.GET("/group/join", api.GroupJoin)
	group.GET("/contacts", api.ContactQuery)
	group.POST("/verify", api.ContactVerify)
	group.POST("/send", api.MessageSend)
	group.GET("/messages", api.MessageQuery)
	group.GET("/quickreply", api.QuickReplyQuery)
	group.POST("/quickreply", api.QuickReplyAdd)
	group.PUT("/quickreply", api.QuickReplyEdit)
	group.DELETE("/quickreply", api.QuickReplyDelete)
	group.GET("/autoreply", api.AutoReplyQuery)
	group.POST("/autoreply", api.AutoReplyAdd)
	group.PUT("/autoreply", api.AutoReplyEdit)
	group.DELETE("/autoreply", api.AutoReplyDelete)

	log.Println(g.Run(*host + ":" + *port))
}