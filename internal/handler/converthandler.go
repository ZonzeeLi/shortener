package handler

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/rest/httpx"
	"shortener/internal/logic"
	"shortener/internal/svc"
	"shortener/internal/types"
)

func ConvertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.ConvertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// 参数规则校验
		if err := validator.New().StructCtx(r.Context(), &req); err != nil {
			logx.Error("validator check failed", logx.Field("err", err.Error()))
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 执行业务逻辑
		l := logic.NewConvertLogic(r.Context(), svcCtx)
		resp, err := l.Convert(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}