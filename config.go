package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis"
)

type redisCfg struct {
	Address  string `json:"address"`
	DB       int    `json:"db"`
}

type bindCfg struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type appCfg struct {
	redisPassword string `json:"-"`
	Api               string   `json:"api_url"`
	CacheTtlSeconds   int      `json:"cache_ttl_seconds"`
	AppTimeoutSeconds int      `json:"app_timeout_seconds"`
	Redis             redisCfg `json:"redis"`
	Bind              bindCfg  `json:"bind"`
}

func (a *appCfg) GetRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     a.Redis.Address,
		Password: a.redisPassword,
		DB:       a.Redis.DB,
	}
}

func (a *appCfg) GetBind() string {
	return fmt.Sprintf("%s:%v", a.Bind.Ip, a.Bind.Port)
}

func (a *appCfg) GetAppTimeout() time.Duration {
	return time.Duration(a.AppTimeoutSeconds) * time.Second
}

func (a *appCfg) GetCacheTtl() time.Duration {
	return time.Duration(a.CacheTtlSeconds) * time.Second
}

func getConfig(filename string, redisPassword string) (*appCfg, error) {
	cfgData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &appCfg{
		redisPassword: redisPassword,
	}
	err = json.Unmarshal(cfgData, cfg)
	return cfg, err
}
