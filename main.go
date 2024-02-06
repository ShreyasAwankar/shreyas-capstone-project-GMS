package main

import (
	"github.com/Capstone/functions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Capstone/docs"
)

// @title Grocery Management System
// @version 1.0.0
// @description Grocery Management API (v1) is a RESTful service that allows you to manage grocery data.
// @description This API provides a set of endpoints to create, retrieve, update, and delete grocery records. It is designed to streamline the processes of managing grocery information in a grocery shop.

// @host us-central1-capstone-412502.cloudfunctions.net
// @schemes https
// @SecurityDefinitions.apikey bearerToken
// @in header
// @name Authorization
func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	grocery := r.Group("api/v1")

	{
		// Define a route
		grocery.POST("/login", functions.Login)
		grocery.GET("/logout", functions.Logout)
		grocery.POST("/item", functions.CreateItem)
		grocery.GET("/grocery", functions.ListGrocery)
		grocery.GET("/item", functions.GetItemById)
		grocery.PUT("/item", functions.UpdateItem)
		grocery.DELETE("/item", functions.DeleteItem)
		grocery.POST("/grocery", functions.CreateBulkUpload)
	}

	r.Run(":8800")
}
