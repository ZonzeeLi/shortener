package logic

import (
	"context"
	"database/sql"
	"errors"
	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	Err404 = errors.New("404")
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 自己写缓存 surl -> lurl
// go-zero自带的缓存 surl -> 数据行

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// 查看短链接，输入 q1mi.com/lusytc -> 重定向到真实的长链接
	// req.ShortUrl = lusytc
	// 1. 根据短链接查询长链接
	// 1.0 布隆过滤器
	// 不存在的短链接直接返回404，不需要后续处理。
	// a. 基于内存版本，缺点：重启后数据丢失，每次重启后都要加载一下已有的短链接（从数据库查）。
	// b. 基于redis版本，go-zero自带的缓存。
	exist, err := l.svcCtx.Filter.ExistsCtx(l.ctx, []byte(req.ShortUrl))
	if err != nil {
		logx.Errorw("BloomFilter.ExistsCtx failed", logx.Field("err", err.Error()))
	}
	if !exist {
		return nil, Err404
	}
	// fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	// 1.1 查询数据库之前可增加缓存层
	// go-zero的缓存支持singleflight，防止缓存击穿
	// 使用singleflight，第一个请求会去查数据库，后续请求会等待第一个请求的结果。
	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{Valid: true, String: req.ShortUrl})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("404")
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed", logx.Field("err", err.Error()))
		return nil, err
	}
	// 2. 返回重定向响应
	return &types.ShowResponse{
		LongUrl: u.Lurl.String,
	}, nil
}
