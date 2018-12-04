package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

func OptionsHandler(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	ctx.Response.Header.Set("Access-Control-Max-Age", "3600")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	ctx.Response.SetStatusCode(200)
}
