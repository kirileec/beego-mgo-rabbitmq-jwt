package main

import (
	"beego-mgo-rabbitmq-jwt/mqHandler"
	_ "beego-mgo-rabbitmq-jwt/routers"
	"beego-mgo-rabbitmq-jwt/utilities/mgodb"
	"beego-mgo-rabbitmq-jwt/utilities/rabbitmq"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/goinggo/tracelog"
)

func main() {
	tracelog.Start(tracelog.LevelTrace)
	tracelog.Started("main", "Initializing Mongo")
	mgodb.Startup("main")
	err := mgodb.Startup("main")
	if err != nil {
		tracelog.CompletedError(err, "main", "initApp")
		os.Exit(1)
	}
	//start listen
	rabbitmq.BeginListenMq(mqHandler.HandleMsg)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	//add ajax Access-Control-Allow-Origin support
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Authorization", "Secure", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:   []string{"Content-Length", "Access-Control-Allow-Origin"},
	}))
	beego.Run()
}
