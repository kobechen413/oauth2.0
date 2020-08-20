package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

var (
	appID     = "9d8a839640521e81a0b9"
	appSecret = "c07aeeb07fa008788f0e46fe748b3123e7a8e30a"
)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/oauth", oauthGithub)

	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("./template/index.tmpl")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	// 利用给定数据渲染模板，并将结果写入w
	url := "https://github.com/login/oauth/authorize?client_id=" + appID
	tmpl.Execute(w, url)
}

func oauthGithub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")

	fmt.Println("authcode", code)

	token, _ := getAuthToken(code)
	userInfo, _ := getUserInfo(token)

	w.Header().Set("content-type", "application/json")

	data, _ := json.Marshal(userInfo)

	w.Write(data)

}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"` // 这个字段没用到
	Scope       string `json:"scope"`      // 这个字段也没用到
}

// 获取令牌
func getAuthToken(autoCode string) (string, error) {
	var url = "https://github.com/login/oauth/access_token?"

	url += "client_id=" + appID + "&client_secret=" + appSecret + "&code=" + autoCode

	fmt.Println("getAuthToken :", url)

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("accept", "application/json")

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer resp.Body.Close()

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		fmt.Println(err)
		return "", err
	}
	return token.AccessToken, nil
}

func getUserInfo(token string) (map[string]interface{}, error) {
	var url = "https://api.github.com/user"

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 将响应的数据写入 userInfo 中，并返回
	var userInfo = make(map[string]interface{})
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	fmt.Println(userInfo)

	return userInfo, nil

}
