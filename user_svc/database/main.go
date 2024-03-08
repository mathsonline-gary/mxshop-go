package main

import (
	"mxshop-go/user_svc/global"
	"mxshop-go/user_svc/model"
)

func main() {
	if err := global.DB.AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}
}
