package main

type AuthReq struct {
	AppID        string `json:"appID"`
	AppKey       string `json:"appKey"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	//	Token        string
	Command string `json:"command"` // 'tail' or 'cat'
	//	UserID       string
	//	Level        string
	//	DateFrom     string
	//	DateTo       string
}
