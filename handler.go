package main

import (
	"apiettiquete/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CachedHandler ...
func CachedHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// CachedUSD ...
func CachedUSD(c *gin.Context) {

	conf, ok := c.MustGet("configuration").(*config.Config)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the configuration"})
		return
	}

	item := cache.Get(conf.URI1)
	if item == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the resource! This should never happen"})
		return
	}

	response := item.Value().([]byte)

	c.Data(200, "application/json; charset=utf-8", response)

}

// CachedEUR ...
func CachedEUR(c *gin.Context) {
	c.String(http.StatusOK, "cached usd")
}
