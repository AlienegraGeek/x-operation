package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"

	"x-operation/internal/svc"
	"x-operation/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type X_operationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewX_operationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *X_operationLogic {
	return &X_operationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *X_operationLogic) X_operation(req *types.Request) (resp *types.Response, err error) {
	// todo: add logic here and delete this line

	return
}

func GetMyUser(accessToken string) (types.UserInfo, error) {
	req, _ := http.NewRequest("GET", "https://api.x.com/2/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	proxyUrl, _ := neturl.Parse("http://127.0.0.1:7890") // 根据你的代理端口改
	transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	httpClient := &http.Client{Transport: transport}
	resp, err := httpClient.Do(req)

	if err != nil {
		return types.UserInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.UserInfo{}, fmt.Errorf("failed to fetch user info: %d", resp.StatusCode)
	}

	var data struct {
		Data types.UserInfo `json:"data"`
	}

	json.NewDecoder(resp.Body).Decode(&data)

	return data.Data, nil
}

func IsFollowing(accessToken, myID, targetID string) (bool, error) {
	// Pro Limit 50 requests / 15 mins PER USER
	// Basic Limit 5 requests / 15 mins PER USER
	// Free Limit ❌
	url := fmt.Sprintf("https://api.x.com/2/users/%s/following", myID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	proxyUrl, _ := neturl.Parse("http://127.0.0.1:7890") // 代理端口
	transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	httpClient := &http.Client{Transport: transport}

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("X API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	for _, u := range result.Data {
		if u.ID == targetID {
			return true, nil
		}
	}
	return false, nil
}
