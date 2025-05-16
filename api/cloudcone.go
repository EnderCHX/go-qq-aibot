package api

import "net/http"

func GetServerInfo(AppSecret, Hash string) {
	req, err := http.NewRequest("GET", "https://api.cloudcone.com/v1/server/info", nil)
}
