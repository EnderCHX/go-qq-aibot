package qqbot

import (
	"fmt"

	"github.com/EnderCHX/go-qq-aibot/ai"
	"github.com/EnderCHX/go-qq-aibot/config"
	"github.com/EnderCHX/go-qq-aibot/search"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func init() {
	c := config.GetConfig()
	deepSeek := ai.DeepSeek{}
	deepSeek.Init(c.DeepSeek.ApiUrl, c.DeepSeek.ApiKey, c.DeepSeek.Model, c.DeepSeek.SysPrompt)

	s := search.NewSearXNG(c.WebSearch.ApiUrl)

	zero.OnMessage().Handle(func(ctx *zero.Ctx) {
		if ctx.Event.IsToMe {

			rcv := ctx.MessageString()

			for rcv[0] == ' ' {
				rcv = rcv[1:]
			}

			if rcv[:len("websearch")] != "websearch" {

				msg, _ := deepSeek.GetMessage(rcv)
				ctx.Send(msg)

				return
			}

			rcv = rcv[len("websearch"):]

			for rcv[0] == ' ' {
				rcv = rcv[1:]
			}

			// fmt.Println(rcv)
			r, err := s.Search(rcv)
			if err != nil {
				ctx.Send("搜索失败")
				return
			}

			msg, _ := deepSeek.GetMessage(
				fmt.Sprintf("这是问题：%s 的联网搜索结果，%s，根据搜索结果回答问题",
					ctx.MessageString(),
					r.ToResultsContent().GetContents().ToString(),
				))
			ctx.Send(msg + " [联网搜索]")
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
