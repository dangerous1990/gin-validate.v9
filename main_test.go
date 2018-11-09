package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCheckMails(t *testing.T) {
	router := setupRouter()
	router.POST("/check", checkMails)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/check", strings.NewReader(`{
		"mails":"1xxx"
	}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error":"Mails邮箱不合法，多个邮箱请以逗号(半角)分隔"}`, w.Body.String())
}

func TestTrueCheckMails(t *testing.T) {
	router := setupRouter()
	router.POST("/check", checkMails)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/check", strings.NewReader(`{
		"mails":"123@163.com,456@163.com"
	}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"form mail are valid!"}`, w.Body.String())
}

type Form struct {
	Mails string `json:"mails" binding:"required,isValidMultiEmails" `
}

func checkMails(c *gin.Context) {
	var form Form
	if err := c.ShouldBindJSON(&form); err == nil {
		fmt.Println(form)
		c.JSON(http.StatusOK, gin.H{"message": "form mail are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
