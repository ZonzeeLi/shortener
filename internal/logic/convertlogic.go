package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	"shortener/pkg/urltool"

	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert short url to long url
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 1. 校验输入的数据
	// 1.1 数据不能空
	// 使用validator库进行校验
	// 1.2 输入的长链接必须是一个能请求通的网址
	if ok := connect.Get(req.LongUrl); !ok {
		return nil, errors.New("无效的链接")
	}
	// 1.3 判断之前是否已经转链过（数据库中是否已存在该长链接）
	// 1.3.1 给长链接生成md5
	md5value := md5.Sum([]byte(req.LongUrl))
	// 1.3.2 拿md5去数据库中查询
	u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5value, Valid: true})
	if err != sqlx.ErrNotFound {
		// 如果存在，且err为nil，说明数据库中已存在该长链接
		if err == nil {
			return nil, fmt.Errorf("该链接已经转过链了，短链为：%s", u.Surl.String)
		}
		logx.Errorw("ShortUrlModel.FindOneByMd5", logx.Field("err", err.Error()))
		return nil, err
	}
	// 1.4 输入的不能是一个短链接（避免循环转链）
	// 输入的是一个完整的url q1mi.com/1d12a?name=q1mi
	basePath, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("url.Parse failed", logx.Field("err", err.Error()), logx.Field("lurl", req.LongUrl))
		return nil, err
	}
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		// 如果存在，且err为nil，说明已经是一个短链接了
		if err == nil {
			return nil, fmt.Errorf("该链接已经是短链接")
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed", logx.Field("err", err.Error()))
		return nil, err
	}
	var short string
	for {

		// 2. 取号 基于MySQL实现的发号器
		// 每来一个转链请求，我们就使用 REPLACE INTO 语句往 sequence 表插入一条数据，然后再使用 SELECT LAST_INSERT_ID() 语句获取刚刚插入的数据的自增 ID。
		seq, err := l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw("Sequence.Next() failed", logx.Field("err", err.Error()))
			return nil, err
		}
		// 3. 号码转短链
		// 3.1 安全性
		// 3.2 短域名黑名单，避免某些特殊词
		short = base62.IntToString(seq)
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
			break // 生成不在黑名单里的短链接就跳出for循环
		}
	}
	// 4. 存储长短链接映射关系
	if _, err := l.svcCtx.ShortUrlModel.Insert(
		l.ctx,
		&model.ShortUrlMap{
			Lurl: sql.NullString{String: req.LongUrl, Valid: true},
			Surl: sql.NullString{String: short, Valid: true},
			Md5:  sql.NullString{String: md5value, Valid: true},
		},
	); err != nil {
		logx.Errorw("ShortUrlModel.Insert failed", logx.Field("err", err.Error()))
	}
	// 将生成的短链接加入到布隆过滤器
	if err := l.svcCtx.Filter.AddCtx(l.ctx, []byte(short)); err != nil {
		logx.Errorw("Filter.AddCtx failed", logx.Field("err", err.Error()))
	}
	// 5. 返回响应
	// 5.1 返回的是 短域名+短链接
	shortUrl := l.svcCtx.Config.ShortDomain + "/" + short
	return &types.ConvertResponse{
		ShortUrl: shortUrl,
	}, nil
}
