package controllers

import (
	"goravel/app/agents"
	"goravel/app/facades"

	"github.com/goravel/framework/contracts/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (r *UserController) Index(ctx http.Context) http.Response {
	conv, err := facades.AI().Agent(&agents.Test{})
	if err != nil {
		panic(err)
	}

	stream, err := conv.Stream("Returns 100 words about Goravel.")
	if err != nil {
		panic(err)
	}

	return stream.HTTPResponse(ctx)

	// return ctx.Response().Success().Json(http.Json{
	// 	"Hello": "Goravel",
	// })
}
