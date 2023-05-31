package svc

import (
	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"
)

type ServiceContext struct {
	Config        config.Config
	ShortUrlModel model.ShortUrlMapModel

	Sequence sequence.Sequence

	ShortUrlBlackList map[string]struct{}

	// bloom filter
	Filter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}
	// 初始化布隆过滤器
	store := redis.MustNewRedis(c.CacheRedis[0].RedisConf, func(r *redis.Redis) {
		r.Type = redis.NodeType
	})
	// 声明一个bitSet，key="bloom_filter"名且bits是1024位。
	filter := bloom.New(store, "bloom_filter", 20*(1<<20))

	// 加载已有的短链接数据
	return &ServiceContext{
		Config:            c,
		ShortUrlModel:     model.NewShortUrlMapModel(conn, c.CacheRedis),
		Sequence:          sequence.NewMySQL(c.Sequence.Dsn),
		ShortUrlBlackList: m,
		Filter:            filter,
	}
}

// 加载已有的短链接数据
func loadDataToBloomFilter() {

}
