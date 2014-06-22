package main

import (
	"fmt"
	"io/ioutil"
)

type UserCreationRequest struct {
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
}

func CreateUser(loginname string, password string) {
	path := fmt.Sprintf("/apps/%s/users", globalConfig.AppId)
	headers := globalConfig.HttpHeaders("application/json")
	req := &UserCreationRequest{loginname, password}
	res := HttpPostJson(path, headers, req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
