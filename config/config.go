package config

import (
	"encoding/json"
	"os"
)

type WebSearch struct {
	ApiUrl string `json:"api_url"`
	Cookie string `json:"cookie"`
}

type QQBot struct {
	WSAddr        string   `json:"ws_addr"`
	Key           string   `json:"key"`
	CommandPrefix string   `json:"command_prefix"`
	NickName      []string `json:"nickname"`
	SuperUsers    []int64  `json:"superusers"`
}

type DeepSeek struct {
	ApiUrl    string `json:"api_url"`
	ApiKey    string `json:"api_key"`
	Model     string `json:"model"`
	SysPrompt string `json:"sys_prompt"`
}

type Config struct {
	QQBot     QQBot     `json:"qqbot"`
	DeepSeek  DeepSeek  `json:"deepseek"`
	WebSearch WebSearch `json:"websearch"`
}

var gbcf *Config

func init() {
	buffer, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	gbcf = &Config{}
	json.Unmarshal(buffer, gbcf)
	// fmt.Println(gbcf)
}

func GetConfig() *Config {
	return gbcf
}
