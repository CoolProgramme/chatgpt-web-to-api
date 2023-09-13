package imitate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/api/chatgpt"
	"github.com/linweiyuan/go-logger/logger"
	"golang.org/x/time/rate"
	"log"
	"sync"
	"time"
)

var loginLimiter = rate.NewLimiter(rate.Every(time.Minute), 1)

func AdminUserAdd(c *gin.Context) {
	var loginInfo []api.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(api.AdminAddUserInfoErrorMessage))
		return
	}
	for _, info := range loginInfo {
		users[info.Username] = info
	}

	go func() {
		for _, info := range loginInfo {
			accessToken, err := TokenGenerate(info)
			if err != nil {
				logger.Error("Token生成失败！")
				continue
			}
			tokensAppend(accessToken)
			logger.Info("Token新增成功！")
		}
	}()

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success!",
		"data":    nil,
	})
}

func AdminTokenAdd(c *gin.Context) {
	var accessTokens []string
	if err := c.ShouldBindJSON(&accessTokens); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(api.AdminAddUserInfoErrorMessage))
		return
	}
	tokensAppendArray(accessTokens)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success!",
		"data":    nil,
	})
}

func AdminTokenGet(c *gin.Context) {
	var loginInfo api.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(api.ParseUserInfoErrorMessage))
		return
	}

	statusCode, errorMessage, accessTokenResponse := chatgpt.GetAccessToken(loginInfo)
	if statusCode != http.StatusOK {
		c.AbortWithStatusJSON(statusCode, api.ReturnMessage(errorMessage))
		return
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(accessTokenResponse), &data)
	if err != nil {
		logger.Error("Json解析错误！")
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success!",
		"data":    data["accessToken"].(string),
	})
}

func AdminTokenCount(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success!",
		"data":    len(tokens),
	})
}

func TokenGenerate(loginInfo api.LoginInfo) (string, error) {
	if err := loginLimiter.Wait(context.Background()); err != nil {
		log.Fatal(err)
	}
	statusCode, errorMessage, accessTokenResponse := chatgpt.GetAccessToken(loginInfo)
	if statusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("%d=%s", statusCode, errorMessage))
		return "", errors.New("login error")
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(accessTokenResponse), &data)
	if err != nil {
		logger.Error("Json解析错误！")
	}
	return data["accessToken"].(string), nil
}

var tokensAppendLock sync.Mutex

func tokensAppendArray(elems []string) {
	tokensAppendLock.Lock()
	defer tokensAppendLock.Unlock()
	tokens = append(tokens, elems...)
}

func tokensAppend(elem string) {
	tokensAppendLock.Lock()
	defer tokensAppendLock.Unlock()
	tokens = append(tokens, elem)
}
