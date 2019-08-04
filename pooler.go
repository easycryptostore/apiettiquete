package main

import (
	"apiettiquete/config"
	"apiettiquete/other"
	"fmt"
	"time"

	"github.com/karlseguin/ccache"
	log "github.com/sirupsen/logrus"
)

// StartPooling ...
func StartPooling(conf *config.Config, cache *ccache.Cache) {

	getToCache(cache, conf.URI1)
	getToCache(cache, conf.URI2)

	ticker := time.NewTicker(time.Duration(conf.PoolIntervalSeconds) * time.Second)
	go func() {
		for range ticker.C {
			getToCache(cache, conf.URI1)
			getToCache(cache, conf.URI2)
			//... add more here
		}
	}()
}

// stores the response to cache
func getToCache(cache *ccache.Cache, url string) {
	fmt.Println("Processing " + url)

	resp, err := other.GetRequest(url)
	if err != nil {
		log.Error(err)
		return
	}

	cache.Set(url, resp, time.Hour*3)

}
