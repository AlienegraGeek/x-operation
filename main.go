package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"x-operation/internal/logic"
	"x-operation/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	RedirectURI = "http://localhost:8080/x/callback"
	Scope       = "tweet.read users.read follows.read offline.access"
	//CodeChallenge = "testchallenge" // 暂时写死，先跑通流程
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	r := gin.Default()
	var verifierStore = map[string]string{}
	// 1. 跳转到 X 登录授权
	r.GET("/x/login", func(c *gin.Context) {
		//state := "xyz" // 最好改成随机，并存起来做CSRF校验
		// 生成随机 state
		state := utils.RandomString(10)

		codeVerifier := utils.GenerateCodeVerifier()
		codeChallenge := utils.GenerateCodeChallenge(codeVerifier)
		// ✅ 纯 plain 模式：直接生成 code_verifier，并作为 code_challenge
		//codeVerifier := utils.RandomString(43) // RFC 要求长度 43~128
		//codeChallenge := codeVerifier          // <-- 就是它！一模一样！
		// 保存 codeVerifier 以便 callback 用
		verifierStore[state] = codeVerifier

		authURL := fmt.Sprintf(
			"https://x.com/i/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=S256",
			os.Getenv("CLIENT_ID"),
			url.QueryEscape(RedirectURI),
			url.QueryEscape(Scope),
			state,
			codeChallenge,
		)
		fmt.Printf("Redirecting to: %s\n", authURL)
		c.Redirect(http.StatusFound, authURL)
	})

	// 2. 获取 code → 换 token
	r.GET("/x/callback", func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")

		codeVerifier, ok := verifierStore[state]
		if !ok {
			c.String(400, "invalid state or missing verifier")
			return
		}

		// ✅ 配置代理
		proxyURL, _ := url.Parse("http://127.0.0.1:7890") // ← 如果你代理是 clash、v2ray，一般就是这个端口
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}

		// ✅ 改用 NewRequest 而不是 PostForm
		form := url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {os.Getenv("CLIENT_ID")},
			"code":          {code},
			"redirect_uri":  {RedirectURI},
			"code_verifier": {codeVerifier},
		}

		req, _ := http.NewRequest("POST", "https://api.x.com/2/oauth2/token", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// ✅ 添加认证头：Basic base64(client_id:client_secret)
		auth := os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET")
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))

		// ✅ 使用带代理的 client 发送
		resp, err := httpClient.Do(req)
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
		//{
		//	"access_token": "YWtsYnFqeURNQ2Jmc2IxMTlOZHp0bmo4V19JcjNJT1dzWm55UGI0cGVoR2l5OjE3NjI0MjQ2ODg5Nzg6MToxOmF0OjE",
		//	"expires_in": 7200,
		//	"refresh_token": "SVROcnZ6TFpCVG5Oc3JpM29aWWxXb0gtWDNxZWt3NHBDVjFBX1p4RUV6bzNuOjE3NjI0MjQ2ODg5Nzg6MToxOnJ0OjE",
		//	"scope": "users.read follows.read tweet.read offline.access",
		//	"token_type": "bearer"
		//}

		c.JSON(200, body)
	})

	r.GET("/x/check_follow", func(c *gin.Context) {
		accessToken := c.Query("token")
		target := c.Query("target") // 如：44196397

		user, err := logic.GetMyUser(accessToken)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ok, err := logic.IsFollowing(accessToken, user.ID, target)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"id":        user.ID,
			"name":      user.Username,
			"target":    target,
			"following": ok,
		})
	})

	r.Run(":8080")
}
