package main

import (
	"apiettiquete/config"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

	var raw []map[string]string
	err := json.Unmarshal(response, &raw)
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, raw)

}

// CachedEUR ...
func CachedEUR(c *gin.Context) {
	c.String(http.StatusOK, "cached usd")
}
