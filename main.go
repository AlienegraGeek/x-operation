package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"x-operation/internal/logic"

	"github.com/gin-gonic/gin"
)

const (
	ClientID      = ""
	RedirectURI   = "http://localhost:8080/x/callback"
	Scope         = "tweet.read users.read follows.read offline.access"
	CodeChallenge = "testchallenge" // 暂时写死，先跑通流程
)

func main() {
	r := gin.Default()

	// 1. 跳转到 X 登录授权
	r.GET("/x/login", func(c *gin.Context) {
		authURL := fmt.Sprintf(
			"https://x.com/i/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=xyz&code_challenge=%s&code_challenge_method=plain",
			ClientID,
			url.QueryEscape(RedirectURI),
			url.QueryEscape(Scope),
			CodeChallenge,
		)
		c.Redirect(http.StatusFound, authURL)
	})

	// 2. 获取 code → 换 token
	r.GET("/x/callback", func(c *gin.Context) {
		code := c.Query("code")

		resp, err := http.PostForm("https://api.twitter.com/2/oauth2/token", url.Values{
			"client_id":     {ClientID},
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {RedirectURI},
			"code_verifier": {CodeChallenge},
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body)
		fmt.Println("===== X OAuth Token Result =====")
		j, _ := json.MarshalIndent(body, "", "  ")
		fmt.Println(string(j))
		fmt.Println("================================")
		// 测试阶段：直接输出到浏览器
		c.JSON(200, body)
	})

	r.GET("/x/check_follow", func(c *gin.Context) {
		accessToken := c.Query("token")
		target := c.Query("target") // 如：44196397

		myID, err := logic.GetMyTwitterID(accessToken)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ok, err := logic.IsFollowing(accessToken, myID, target)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"user_id":   myID,
			"target":    target,
			"following": ok,
		})
	})

	r.Run(":8080")
}
