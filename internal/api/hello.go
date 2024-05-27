package api

import (
	"gin-starter/internal/global"
	"gin-starter/internal/model/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Hello(c *gin.Context) {
	time.Sleep(time.Second * 2)
	c.JSON(http.StatusOK, gin.H{
		"now": time.Now().Format(time.DateTime),
	})
}

func HelloPost(c *gin.Context) {
	type RegisterRequest struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Age      uint8  `json:"age" binding:"gte=1,lte=120"`
	}
	var r RegisterRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.ValidatorError(c, err)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"title": "提交正确",
		})
	}
}

type Page struct {
	Page int `form:"page" binding:"required,gte=1"`
	Size int `form:"size" binding:"required,gte=1"`
}

type FatherReq struct {
	FatherName string `json:"name"`
	FatherAge  int    `json:"age"`
	SonName    string `json:"son_name"`
	SonAge     int    `json:"son_age"`
}

func TestGorm(c *gin.Context) {
	var fathersWithSons []FatherReq
	var pager Page
	if err := c.ShouldBind(&pager); err != nil {
		response.ValidatorError(c, err)
		return
	}
	page := pager.Page
	size := pager.Size
	offset := (page - 1) * size

	global.DB.Table("father").Select("father.name, father.age, son.age as son_age, son.name as son_name").
		Joins("INNER JOIN son ON father.id = son.father_id").
		Where("son.age > ?", 10).
		Limit(size).
		Offset(offset).
		Scan(&fathersWithSons)

	c.JSON(200, fathersWithSons)
}
