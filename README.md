# Gin Auth middleware
[![Build Status](https://travis-ci.org/nikonm/go-gin-auth-middleware.svg?branch=master)](https://travis-ci.org/nikonm/go-gin-auth-middleware)

Library implemented auth middleware for gin application

With security providers:
 - DB
 - FreeIPA
 
 ### Usage 
 ```golang
 package main
 
 import (
 	"github.com/gin-gonic/gin"
 	"github.com/nikonm/go-gin-auth-middleware"
 	"github.com/nikonm/go-gin-auth-middleware/user"
 	"log"
 	"net/http"
 )
 
 func main()  {
 
 
 	secOptions := security.Options{
 		Secret:     "SomeSecretKey",
 		TokenExp:   3600, //time.Duration
 		HeaderName: "X-AUTH-TOKEN",
 		Adapters: map[string]map[string]interface{}{
 			"db": {
 				"driver": "postgres",
 				"connection": "postgres://user:pass@localhost/test-database",
 				"sql": "SELECT {select_columns} FROM users where username=$1 and password=$2",
 				"source_target_fields": map[string]interface{}{ //Mapping source fields from adapter and target fields from User model (github.com/nikonm/go-gin-auth-middleware/user)
 					"id": "id",
 					"username": "user",
 					"email": "email",
 				},
 			},
 			"ipa": {
 				"host": "ipa.test.local",
 				"timeout": "1m",
 				"secured": false,
 				"source_target_fields": map[string]interface{}{ //Same as db adapter
 					"id": "id",
 					"login": "login",
 					"username": "user",
 				},
 			},
 
 		},
 	}
 	err := security.Init(secOptions)
 	if err != nil {
 		log.Fatal(err)
 	}
 
 	r := gin.Default()
 	r.POST("/login", func(c *gin.Context) {
 		login := user.LoginDTO{}
 		err := c.BindJSON(&login)
 		if err != nil {
 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
 			return
 		}
 
 		token, err := security.Security.Login(c, login)
 		if err != nil {
 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
 			return
 		}
 		c.JSON(200, gin.H{
 			"token": token,
 		})
 
 	})
 	r.GET("/secured", func(c *gin.Context) {
 		u, _ := security.Security.GetUser(c)
 
 		c.JSON(200, gin.H{
 			"user": u,
 		})
 	}).Use(security.Security.Middleware)
 	r.Run()
 }

 
 ```