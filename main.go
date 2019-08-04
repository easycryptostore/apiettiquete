package main

import (
	"apiettiquete/config"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/karlseguin/ccache"
)

var cache = ccache.New(ccache.Configure().MaxSize(1000).ItemsToPrune(500))

var conf *config.Config
var router *gin.Engine

// load the configuration
func init() {

	env := flag.String("env", "", "production or development")
	flag.Parse()

	switch *env {
	case "production":
		if err := godotenv.Load("./production.env"); err != nil {
			log.Fatal("no production.env file found. Exiting...")
		}
		break
	case "development":
		if err := godotenv.Load("./development.env"); err != nil {
			log.Fatal("no development.env file found. Exiting...")
		}
		break
	default:
		log.Fatal("please specify -env=production (or development)")
	}

	conf = config.New()
}

func main() {

	log.Println("APIettiquete")
	log.Printf("running on http://%s:%d", conf.ServerHostname, conf.ServerPort)

	//-------- LOGIC --------
	err := StartPooling(conf, cache)
	if err != nil {
		log.Fatal(err)
	}

	//-------- API --------
	if conf.DebugMode {
		gin.SetMode(gin.DebugMode)
		log.SetLevel(log.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.WarnLevel)
	}

	router = gin.New()
	router.Use(gin.Logger())
	router.Use(corsMiddleware())
	router.Use(configurationMiddleware(conf)) // NEEDED ?
	router.Use(gin.Recovery())

	api := router.Group("/api")

	api.GET("/USD", CachedUSD)
	api.GET("/EUR", CachedEUR)

	server := &http.Server{
		Addr:           conf.ServerHostname + ":" + strconv.Itoa(conf.ServerPort),
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1Mb
	}
	server.SetKeepAlivesEnabled(true)

	// Serve'em
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("initiated server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
	}
	log.Println("server exiting. bye!")

}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// configurationMiddleware will add the configuration to the context
func configurationMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("configuration", config)
		c.Next()
	}
}
