package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

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
	// todo: add your logic here and delete this line

	return
}

func GetMyTwitterID(accessToken string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.twitter.com/2/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	proxyUrl, _ := url.Parse("http://127.0.0.1:7890") // 根据你的代理端口改
	transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	httpClient := &http.Client{Transport: transport}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Data struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&data)

	return data.Data.ID, nil
}

func IsFollowing(accessToken, myID, targetID string) (bool, error) {
	url1 := fmt.Sprintf("https://api.twitter.com/2/users/%s/following", myID)

	req, _ := http.NewRequest("GET", url1, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	proxyUrl, _ := url.Parse("http://127.0.0.1:7890") // 根据你的代理端口改
	transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	httpClient := &http.Client{Transport: transport}
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

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
