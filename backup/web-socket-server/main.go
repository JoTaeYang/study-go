package main

import (
	"log"
	"net/http"

	"github.com/JoTaeYang/study-go/library/fsyaml"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

const (
	Default_Path      = "./"
	Default_File_Name = "config.yaml"
)

var (
	cfg    *fsyaml.Config
	engine *gin.Engine
)

func CORSM() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Method", "GET, DELETE, POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	}
}

func InitConfig() error {
	cfg = &fsyaml.Config{}
	path := Default_Path + Default_File_Name
	err := fsyaml.Init(cfg, path)
	if nil != err {
		log.Println(err)
		return err
	}
	return nil
}

func InitGin() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	pprof.Register(r)

	r.Use(CORSM())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	AcceptRouter := r.Group("yang/ws/accept")
	{
		AcceptRouter.GET("")
	}

	return r
}

func main() {

	// init config
	err := InitConfig()
	if nil != err {
		return
	}
	// gin start
	engine = InitGin()

	engine.Run(cfg.Env.Port)
}
