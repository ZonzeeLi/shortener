# 短连接项目

## 搭建项目的骨架

1. 建库建表

新建发号器表

```sql
CREATE TABLE `sequence` (
                            `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                            `stub` varchar(1) NOT NULL,
                            `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `idx_uniq_stub` (`stub`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT = '序号表';
```

新建长链接短链接映射表：

```sql
CREATE TABLE `short_url_map` (
                                 `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
                                 `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 `create_by` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '创建者',
                                 `is_del` tinyint UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否删除：0正常1删除',

                                 `lurl` varchar(2048) DEFAULT NULL COMMENT '长链接',
                                 `md5` char(32) DEFAULT NULL COMMENT '长链接MD5',
                                 `surl` varchar(11) DEFAULT NULL COMMENT '短链接',
                                 PRIMARY KEY (`id`),
                                 INDEX(`is_del`),
                                 UNIQUE(`md5`),
                                 UNIQUE(`surl`)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT = '长短链映射表';
```

2. 搭建go-zero框架的骨架

2.1 编写`api`文件

```api
type ConvertRequest {
    LongUrl string `json:"longUrl"`
}

type ConvertResponse {
    ShortUrl string `json:"shortUrl"`
}

type ShowRequest {
    ShortUrl string `json:"shortUrl"`
}

type ShowResponse {
    LongUrl string `json:"longUrl"`
}

service shortener-api {

    @handler ConvertHandler
    post /convert(ConvertRequest) returns(ConvertResponse)

    // q1mi.cn/lycsa1
    @handler ShowHandler
    get /:shortUrl(ShowRequest) returns(ShowResponse)

}
```

2.2 根据`api`文件生成`go`文件

```bash
goctl api go -api shortener.api -dir .
```

2.3 根据`sql`文件生成`model`层

```bash
goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/short" -table="sequence"  -dir="./model"

goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/short" -table="short_url_map"  -dir="./model"
```

2.4 下载项目依赖

```bash
go mod tidy
```

3. 运行项目

```bash
go run shortener.go
```

4. 修改配置结构体和配置文件

## 参数校验

1. go-zero使用validator进行参数校验

下载依赖：
```bash
go get github.com/go-playground/validator/v10
```

导入依赖：
```bash
import "github.com/go-playground/validator/v10"
```

在api中为结构体添加validate tag，并添加校验规则

## 查看短链接

### 缓存版

有两种方式，

1. 使用自己实现的缓存，surl -> lurl，能够节省缓存空间，缓存数据量小
2. 使用go-zero自带的缓存，surl -> 数据行，代码量少，开发量小

这里使用第二种方案：
1. 添加缓存配置
    - 配置文件
    - 配置config结构体
2. 删除旧model层代码
    - 删除 shorturlmapmodel.go
3. 生成新的model层代码
    - goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/short" -table="short_url_map"  -dir="./model" -c
4. 修改svccontext层代码

## 项目如何扩展？

1. 如何支持自定义短链？

维护一个已经使用的序号，后续生成序号时，判断是否已经被分配。没有被分配，则使用该序号，否则使用发号器生成的序号。

2. 如何让短链支持过期时间？

每个链接映射额外记录一个过期时间字段，到期后将该映射记录删除。

关于删除的策略有以下几种：
   
- 延迟删除：每次请求时判断是否过期，如果过期则删除。
  - 实现简单，性能损失小。
  - 存储空间的利用效率低，已经过期得数据可能永远不会被删除。
- 定时删除：创建记录时根据过期时间设置定时器。
  - 过期数据能被及时删除，存储空间的利用率高。
  - 占用内存大，性能差。
- 轮询删除：通过异步脚本在业务低峰期周期性扫表清理过期数据
  - 兼顾效率和磁盘利用率。
  - 
3. 如何提高吞吐量？

整个系统分为生成短链（写）和访问短链（读）两部分

- 水平拓展多节点，根据序号分片。

4. 延迟优化

整个系统分为生成短链（写）和访问短链（读）两部分

- 存储层
  - 数据结构简单可以直接改用kv存储
  - 对存储节点进行分片
- 缓存层
  - 增加缓存层，本地缓存-->redis缓存
  - 使用布隆过滤器判断长链接映射是否已存在，判断短链接是否有效
- 网络
  - 基于地理位置就近访问数据节点





