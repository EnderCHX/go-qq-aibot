package qqbot

import (
	"github.com/EnderCHX/go-qq-aibot/ai"
	"github.com/EnderCHX/go-qq-aibot/config"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func init() {
	c := config.GetConfig()
	deepSeek := ai.DeepSeek{}
	deepSeek.Init(c.DeepSeek.ApiUrl, c.DeepSeek.ApiKey, c.DeepSeek.Model, c.DeepSeek.SysPrompt)

	zero.OnMessage().Handle(func(ctx *zero.Ctx) {
		if ctx.Event.IsToMe {
			msg, _ := deepSeek.GetMessage(ctx.MessageString())
			ctx.Send(msg)
		}
	})

	zero.RunAndBlock(&zero.Config{
		NickName:      c.QQBot.NickName,
		CommandPrefix: c.QQBot.CommandPrefix,
		SuperUsers:    c.QQBot.SuperUsers,
		Driver: []zero.Driver{
			driver.NewWebSocketClient(c.QQBot.WSAddr, c.QQBot.Key),
		},
	}, nil)
}
