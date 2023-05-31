package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	ShortUrlDB ShortUrlDB

	Sequence struct {
		Dsn string
	}

	BaseString string

	ShortUrlBlackList []string
	ShortDomain       string

	CacheRedis cache.CacheConf // redis缓存
}

type ShortUrlDB struct {
	DSN string
}
