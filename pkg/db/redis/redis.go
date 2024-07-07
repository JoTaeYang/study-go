package redis

import (
	"errors"
	"log"

	redis "github.com/redis/go-redis/v9"
)

var (
	// WriteRedisClient ...
	WriteRedisClient *redis.ClusterClient
	// ReadRedisClient ...
	ReadRedisClient *redis.ClusterClient

	// pub
	PubRedisClient *redis.ClusterClient

	// sub
	SubRedisClient *redis.ClusterClient
)

var (
	appName = "macovill.develop"
)

type cfgig struct {
	PAddr   []string `yaml:"primary_addr"` // 클러스터 모드 대비해서 Slice 로 주소 설정
	AppName string   `yaml:"app_name"`     //
	Mode    string   `yaml:"mode"`
	PubSub  bool     `yaml:"pub_sub"`
}

/*
레디스 환경 세팅
*/
func Init(cfg cfgig) error {

	if cfg.AppName != "" {
		appName = cfg.AppName
	}

	WriteRedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cfg.PAddr,
		MinIdleConns: 50,
		PoolSize:     100,
	})
	if WriteRedisClient == nil {
		return errors.New("redis connect fail")
	}

	ReadRedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cfg.PAddr,
		PoolSize: 100,
	})
	if ReadRedisClient == nil {
		return errors.New("redis connect fail")
	}

	if cfg.PubSub {
		PubRedisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: cfg.PAddr,
		})
		if PubRedisClient == nil {
			return errors.New("redis cluster fail")
		}
		log.Println("redis pubsub ready")
	}

	log.Println("redis connect success")
	return nil
}
