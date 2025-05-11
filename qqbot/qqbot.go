package qqbot

import (
	"errors"
	"fmt"
	"github.com/EnderCHX/go-qq-aibot/ai/chat"
	"github.com/EnderCHX/go-qq-aibot/config"
	df "github.com/EnderCHX/go-qq-aibot/delta_force"
	armed "github.com/EnderCHX/go-qq-aibot/delta_force/armed_force"
	item "github.com/EnderCHX/go-qq-aibot/delta_force/items"
	"github.com/EnderCHX/go-qq-aibot/delta_force/maps"
	"github.com/EnderCHX/go-qq-aibot/search"
	"github.com/bytedance/sonic"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	c := config.GetConfig()
	deepSeek := chat.DeepSeek{}
	deepSeek.Init(c.DeepSeek.ApiUrl, c.DeepSeek.ApiKey, c.DeepSeek.Model, c.DeepSeek.SysPrompt)

	s := search.NewSearXNG(c.WebSearch.ApiUrl)

	zero.OnCommand("状态").Handle(func(ctx *zero.Ctx) {
		ctx.Send(fmt.Sprintf("状态：%s", time.Now().Format("2006-01-02 15:04:05")))
	})

	zero.OnMessage().Handle(func(ctx *zero.Ctx) {
		if ctx.Event.IsToMe {

			rcv := ctx.MessageString()

			for rcv[0] == ' ' {
				rcv = rcv[1:]
			}

			if len(rcv) <= len("websearch") {
				msg, _ := deepSeek.GetMessage(rcv)
				ctx.Send(msg)
				return
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

	zero.OnMessage().Handle(func(ctx *zero.Ctx) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Send("出现错误：" + fmt.Sprint(err))
			}
		}()

		msg := ctx.MessageString()

		if match, err := regexp.MatchString("^df*", msg); err == nil && match {
			msg = msg[len("df"):]
		} else if match, err := regexp.MatchString("^△*", msg); err == nil && match {
			msg = msg[len("△"):]
		} else if match, err := regexp.MatchString("^三角洲*", msg); err == nil && match {
			msg = msg[len("三角洲"):]
		} else {
			return
		}

		if match, err := regexp.MatchString("^帮助$", msg); err == nil && match {
			//"命令前缀为\"df\"、\"△\"、\"三角洲\"\n用法: 命令前缀加以下命令：\n1. 设置ck\n2. 密码\n3. 近日收益\n 4.战绩"
			ctx.Send([]message.Segment{
				message.Text("命令前缀为\"df\"、\"△\"、\"三角洲\"\n"),
				message.Text("用法: 命令前缀加以下命令：\n"),
				message.Text("1. 设置ck\n"),
				message.Text("2. ck帮助"),
				message.Text("3. 密码\n"),
				message.Text("4. 近日收益\n"),
				message.Text("5. 战绩\n"),
				message.Text("6. 大红\n"),
				message.Text("7. 登录记录(page int)\n"),
				message.Text("8. 道具(page int)\n"),
				message.Text("9. 哈夫币(page int)\n"),
				message.Text("10. 战局(page int)\n"),
			})
		}

		if match, err := regexp.MatchString("^设置ck*", msg); err == nil && match {
			ck := msg[len("设置ck"):]
			ck = strings.TrimSpace(ck)
			err := df.SaveCookie(ck, fmt.Sprint(ctx.Event.UserID))
			if err != nil {
				ctx.SendChain(
					message.Reply(ctx.Event.MessageID),
					message.Text("失败: "+err.Error()),
				)
				return
			}
			ctx.Send("设置成功")
		}

		if match, err := regexp.MatchString("^ck帮助$", msg); err == nil && match {
			ctx.SendChain(
				message.Reply(ctx.Event.MessageID),
				message.Text("浏览器打开 https://df.qq.com/cp/a20241230webmp/index.html, 登录账号\n"),
				message.Text("登录后按f12打开开发者选项, 找到如下请求，复制cookie"),
				message.Image("https://img.020928.xyz/img/df_ck_help.png"),
			)
		}

		if match, err := regexp.MatchString("^大红$", msg); err == nil && match {
			data, err := df.GetAchievements(fmt.Sprint(ctx.Event.UserID))
			if err != nil {
				ctx.SendChain(
					message.Reply(ctx.Event.MessageID),
					message.Text("失败: "+err.Error()),
				)
				return
			}

			jData, _ := sonic.Get(data, "jData", "data", "data", "solDetail")
			redTotalMoney, _ := jData.Get("redTotalMoney").Int64()
			redTotalCount, _ := jData.Get("redTotalCount").Int64()
			redCollectionDetail := jData.Get("redCollectionDetail")
			redCollectionDetailLen := func() int {
				redCollectionDetailArr, _ := redCollectionDetail.Array()
				return len(redCollectionDetailArr)
			}()

			msgChain := []message.Segment{
				message.Reply(ctx.Event.MessageID),
				message.Text(fmt.Sprintf("收藏大红价值: %v\n收藏大红数量: %v\n", redTotalMoney, redTotalCount)),
			}

			for i := 0; i < redCollectionDetailLen; i++ {
				objectID, _ := redCollectionDetail.Index(i).Get("objectID").Int64()
				count, _ := redCollectionDetail.Index(i).Get("count").Int64()
				price, _ := redCollectionDetail.Index(i).Get("price").Int64()
				m := message.Text(fmt.Sprintf("%v(%v): %v\n",
					item.ItemMap[int(objectID)].ObjectName,
					price,
					count,
				))
				msgChain = append(msgChain, m)
			}

			ctx.Send(msgChain)
		}

		if match, err := regexp.MatchString("^战绩$", msg); err == nil && match {
			//achieve_data, err1 := df.GetAchievements(fmt.Sprint(ctx.Event.UserID))
			kd_data, err2 := df.GetKd(fmt.Sprint(ctx.Event.UserID))
			//if err1 != nil {
			//	ctx.SendChain(
			//		message.Reply(ctx.Event.MessageID),
			//		message.Text("失败: "+err1.Error()),
			//	)
			//	return
			//}
			if err2 != nil {
				ctx.SendChain(
					message.Reply(ctx.Event.MessageID),
					message.Text("失败: "+err2.Error()),
				)
				return
			}

			//jData, _ := sonic.Get(achieve_data, "jData", "data", "data", "solDetail")
			//redTotalMoney, _ := jData.Get("redTotalMoney").Int64()
			//redTotalCount, _ := jData.Get("redTotalCount").Int64()
			//redCollectionDetail, _ := jData.Get("redCollectionDetail").Array()

			jData2, _ := sonic.Get(kd_data, "jData")
			picurl, _ := jData2.Get("userData").Get("picurl").String()

			if strings.Contains(picurl, "http") {
				picurl, _ = url.QueryUnescape(picurl)
			} else {
				picurl = "https://playerhub.df.qq.com/playerhub/60004/object/" + picurl + ".png"
			}

			charac_name, _ := jData2.Get("userData").Get("charac_name").String()
			careerData := jData2.Get("careerData")
			rankpoint, _ := careerData.Get("rankpoint").Int64()
			//tdmrankpoint, _ := careerData.Get("tdmrankpoint").Int64()
			soltotalfght, _ := careerData.Get("soltotalfght").Int64()
			solttotalescape, _ := careerData.Get("solttotalescape").Int64()
			solduration, _ := careerData.Get("solduration").Int64()
			soltotalkill, _ := careerData.Get("soltotalkill").Int64()
			solescaperatio, _ := careerData.Get("solescaperatio").String()
			//avgkillperminute, _ := careerData.Get("avgkillperminute").Int64()
			//tdmduration, _ := careerData.Get("tdmduration").Int64()
			//tdmsuccessratio, _ := careerData.Get("tdmsuccessratio").String()
			//tdmtotalfight, _ := careerData.Get("tdmtotalfight").Int64()
			//tdmtotalkill, _ := careerData.Get("tdmtotalkill").Int64()
			ctx.SendChain(
				message.Reply(ctx.Event.MessageID),
				message.Image(picurl),
				message.Text(fmt.Sprintf("角色: %v\n", charac_name)),
				//message.Text(fmt.Sprintf("收藏大红价值: %v\n", redTotalMoney)),
				//message.Text(fmt.Sprintf("收藏大红数量: %v\n", redTotalCount)),
				message.Text(fmt.Sprintf("烽火排位分数: %v\n", rankpoint)),
				message.Text(fmt.Sprintf("烽火总对局: %v\n", soltotalfght)),
				message.Text(fmt.Sprintf("烽火总撤离: %v\n", solttotalescape)),
				message.Text(fmt.Sprintf("烽火总撤离率: %v\n", solescaperatio)),
				message.Text(fmt.Sprintf("烽火总时长: %vh%vmin\n", solduration/60/60, (solduration-(solduration/60/60)*60*60)/60)),
				message.Text(fmt.Sprintf("烽火击杀干员: %v\n", soltotalkill)),
			)
		}

		if match, err := regexp.MatchString("^密码$", msg); err == nil && match {
			data, err := df.GetPassword(fmt.Sprint(ctx.Event.UserID))
			if err != nil {
				ctx.SendChain(
					message.Reply(ctx.Event.MessageID),
					message.Text("失败: "+err.Error()),
				)
				return
			}
			pass, _ := sonic.Get(data, "jData", "data", "data", "content", "secretDay", "data", 0, "desc")
			fmt.Println(pass)
			password, _ := pass.String()
			ctx.SendChain(
				message.Reply(ctx.Event.MessageID),
				message.Text(fmt.Sprintf("密码: \n%v;", password)),
			)
		}

		if match, err := regexp.MatchString("^近日收益$", msg); err == nil && match {
			data, err := df.GetRecent(fmt.Sprint(ctx.Event.UserID))
			if err != nil {
				ctx.SendChain(
					message.Reply(ctx.Event.MessageID),
					message.Text("失败: "+err.Error()),
				)
				return
			}
			solDetail, err := sonic.Get(data, "jData", "data", "data", "solDetail")
			recentGain, _ := solDetail.Get("recentGain").String()
			recentGainDate, _ := solDetail.Get("recentGainDate").String()
			item1Id, _ := solDetail.Get("userCollectionTop").Get("list").Index(0).Get("objectID").Int64()
			item1Count, _ := solDetail.Get("userCollectionTop").Get("list").Index(0).Get("count").Int64()
			item1Price, _ := solDetail.Get("userCollectionTop").Get("list").Index(0).Get("price").Int64()
			item1Name := item.ItemMap[int(item1Id)].ObjectName
			item2Id, _ := solDetail.Get("userCollectionTop").Get("list").Index(1).Get("objectID").Int64()
			item2Count, _ := solDetail.Get("userCollectionTop").Get("list").Index(1).Get("count").Int64()
			item2Price, _ := solDetail.Get("userCollectionTop").Get("list").Index(1).Get("price").Int64()
			item2Name := item.ItemMap[int(item2Id)].ObjectName
			item3Id, _ := solDetail.Get("userCollectionTop").Get("list").Index(2).Get("objectID").Int64()
			item3Count, _ := solDetail.Get("userCollectionTop").Get("list").Index(2).Get("count").Int64()
			item3Price, _ := solDetail.Get("userCollectionTop").Get("list").Index(2).Get("price").Int64()
			item3Name := item.ItemMap[int(item3Id)].ObjectName

			ctx.SendChain(
				message.Reply(ctx.Event.MessageID),
				message.Text(fmt.Sprintf("近日收益: %v\n收益日期: %v\n收集品: \n1. %v(%v): %v\n2. %v(%v): %v\n3. %v(%v): %v", recentGain, recentGainDate,
					item1Name, item1Price, item1Count,
					item2Name, item2Price, item2Count,
					item3Name, item3Price, item3Count,
				)),
			)
		}

		if matchDefault, match, err := func() (bool, bool, error) {
			matchDefault, err1 := regexp.MatchString(`^登录记录$`, msg)
			match, err2 := regexp.MatchString(`^登录记录\d+$`, msg)
			if err1 != nil || err2 != nil {
				return matchDefault, match, errors.New("error: " + err1.Error() + err2.Error())
			}
			return matchDefault, match, nil
		}(); err == nil && (match || matchDefault) {
			var page int
			if matchDefault {
				page = 1
			}
			if match {
				page, _ = strconv.Atoi(msg[len("登录记录"):])
			}

			startIndex := func() int {
				if page%10 != 0 {
					return (page%10 - 1) * 5
				} else {
					return 9 * 5
				}
			}()
			endIndex := func() int {
				if page%10 != 0 {
					return (page % 10) * 5
				} else {
					return 10 * 5
				}
			}()

			data, err := df.GetLoginPropsMoneyBattle(fmt.Sprint(ctx.Event.UserID), page, 1)

			if err != nil {
				ctx.SendChain(
					message.Reply(ctx.Event.MessageID),
					message.Text("失败: "+err.Error()),
				)
				return
			}
			jData, _ := sonic.Get(data, "jData", "data")
			vRoleName, _ := jData.Get("vRoleName").String()
			loginDay, _ := jData.Get("loginDay").String()
			Level, _ := jData.Get("Level").Int64()
			LoginArr := jData.Get("LoginArr")
			LoginArrLen := func() int {
				iA, _ := LoginArr.Array()
				return len(iA)
			}()

			msgChain := []message.Segment{
				message.Reply(ctx.Event.MessageID),
				message.Text(fmt.Sprintf("角色: %v 登录天数: %v 等级: %v\n", vRoleName, loginDay, Level)),
			}

			for i := startIndex; i < endIndex && i < LoginArrLen; i++ {

				indtEventTime, _ := LoginArr.Index(i).Get("indtEventTime").String()
				outdtEventTime, _ := LoginArr.Index(i).Get("outdtEventTime").String()
				vClientIP, _ := LoginArr.Index(i).Get("vClientIP").String()
				SystemHardware, _ := LoginArr.Index(i).Get("SystemHardware").String()

				m := message.Text(fmt.Sprintf("设备%v, IP %v, 上线时间%v, 下线时间%v\n",
					SystemHardware,
					vClientIP,
					indtEventTime,
					outdtEventTime,
				))
				msgChain = append(msgChain, m)
			}
			ctx.Send(msgChain)
		}

		if matchDefault, match, err := func() (bool, bool, error) {
			matchDefault, err1 := regexp.MatchString(`^道具$`, msg)
			match, err2 := regexp.MatchString(`^道具\d+$`, msg)
			if err1 != nil || err2 != nil {
				return matchDefault, match, errors.New("error: " + err1.Error() + err2.Error())
			}
			return matchDefault, match, nil
		}(); err == nil && (match || matchDefault) {
			var page int
			if matchDefault {
				page = 1
			}
			if match {
				page, _ = strconv.Atoi(msg[len("道具"):])
			}

			startIndex := func() int {
				if page%10 != 0 {
					return (page%10 - 1) * 5
				} else {
					return 9 * 5
				}
			}()
			endIndex := func() int {
				if page%10 != 0 {
					return (page % 10) * 5
				} else {
					return 10 * 5
				}
			}()

			data, err := df.GetLoginPropsMoneyBattle(fmt.Sprint(ctx.Event.UserID), page, 2)

			if err != nil {
				ctx.SendChain(
					message.Reply(ctx.Event.MessageID),
					message.Text("失败: "+err.Error()),
				)
				return
			}
			jData, _ := sonic.Get(data, "jData", "data")
			itemArr := jData.Get("itemArr")
			itemArrLen := func() int {
				iA, _ := itemArr.Array()
				return len(iA)
			}()

			msgChain := []message.Segment{
				message.Reply(ctx.Event.MessageID),
			}

			for i := startIndex; i < endIndex && i < itemArrLen; i++ {
				dtEventTime, _ := itemArr.Index(i).Get("dtEventTime").String()
				Name, _ := itemArr.Index(i).Get("Name").String()
				AfterCount, _ := itemArr.Index(i).Get("AfterCount").Int64()
				AddOrReduce, _ := itemArr.Index(i).Get("AddOrReduce").Int64()
				Reason, _ := itemArr.Index(i).Get("Reason").String()
				Reason, _ = url.QueryUnescape(Reason)
				m := message.Text(fmt.Sprintf("在%v %v %v %v 剩余%v\n",
					dtEventTime,
					Reason,
					Name,
					AddOrReduce,
					AfterCount,
				))
				msgChain = append(msgChain, m)
			}
			ctx.Send(msgChain)
		}

		if matchDefault, match, err := func() (bool, bool, error) {
			matchDefault, err1 := regexp.MatchString(`^哈夫币$`, msg)
			match, err2 := regexp.MatchString(`^哈夫币\d+$`, msg)
			if err1 != nil || err2 != nil {
				return matchDefault, match, errors.New("error: " + err1.Error() + err2.Error())
			}
			return matchDefault, match, nil
		}(); err == nil && (match || matchDefault) {
			var page int
			if matchDefault {
				page = 1
			}
			if match {
				page, _ = strconv.Atoi(msg[len("哈夫币"):])
			}

			startIndex := func() int {
				if page%10 != 0 {
					return (page%10 - 1) * 5
				} else {
					return 9 * 5
				}
			}()
			endIndex := func() int {
				if page%10 != 0 {
					return (page % 10) * 5
				} else {
					return 10 * 5
				}
			}()

			data, err := df.GetLoginPropsMoneyBattle(fmt.Sprint(ctx.Event.UserID), page, 3)

			if err != nil {
				ctx.SendChain(
					message.Reply(ctx.Event.MessageID),
					message.Text("失败: "+err.Error()),
				)
				return
			}
			jData, _ := sonic.Get(data, "jData", "data")
			totalMoney, _ := jData.Get("0").Get("totalMoney").Int64()
			iMoneyArr := jData.Get("iMoneyArr")
			iMoneyArrLen := func() int {
				iA, _ := iMoneyArr.Array()
				return len(iA)
			}()

			msgChain := []message.Segment{
				message.Reply(ctx.Event.MessageID),
				message.Text(fmt.Sprintf("总哈夫币: %v\n", totalMoney)),
			}

			for i := startIndex; i < endIndex && i < iMoneyArrLen; i++ {
				dtEventTime, _ := iMoneyArr.Index(i).Get("dtEventTime").String()
				leftMoney, _ := iMoneyArr.Index(i).Get("leftMoney").Int64()
				AddOrReduce, _ := iMoneyArr.Index(i).Get("AddOrReduce").String()
				Reason, _ := iMoneyArr.Index(i).Get("Reason").String()
				Reason, _ = url.QueryUnescape(Reason)
				m := message.Text(fmt.Sprintf("在%v %v %v 剩余%v\n", dtEventTime, Reason, AddOrReduce, leftMoney))
				msgChain = append(msgChain, m)
			}
			ctx.Send(msgChain)
		}

		if matchDefault, match, err := func() (bool, bool, error) {
			matchDefault, err1 := regexp.MatchString(`^战局$`, msg)
			match, err2 := regexp.MatchString(`^战局\d+$`, msg)
			if err1 != nil || err2 != nil {
				return matchDefault, match, errors.New("error: " + err1.Error() + err2.Error())
			}
			return matchDefault, match, nil
		}(); err == nil && (match || matchDefault) {
			var page int
			if matchDefault {
				page = 1
			}
			if match {
				page, _ = strconv.Atoi(msg[len("战局"):])
			}

			startIndex := func() int {
				if page%10 != 0 {
					return (page%10 - 1) * 5
				} else {
					return 9 * 5
				}
			}()
			endIndex := func() int {
				if page%10 != 0 {
					return (page % 10) * 5
				} else {
					return 10 * 5
				}
			}()

			data, err := df.GetLoginPropsMoneyBattle(fmt.Sprint(ctx.Event.UserID), page, 4)

			if err != nil {
				ctx.SendChain(
					message.Reply(ctx.Event.MessageID),
					message.Text("失败: "+err.Error()),
				)
				return
			}
			jData, _ := sonic.Get(data, "jData")
			battleData := jData.Get("data")

			msgChan := []message.Segment{
				message.Reply(ctx.Event.MessageID),
			}

			for i := startIndex; i < endIndex && i < func() int { battleDataArr, _ := battleData.Array(); return len(battleDataArr) }(); i++ {
				battle := battleData.Index(i)
				MapId, _ := battle.Get("MapId").Int64()
				EscapeFailReason, _ := battle.Get("EscapeFailReason").Int64()
				FinalPrice, _ := battle.Get("FinalPrice").Int64()
				//KeyChainCarryOutPrice, _ := battle.Get("KeyChainCarryOutPrice").Int64()
				//CarryoutSafeBoxPrice, _ := battle.Get("CarryoutSafeBoxPrice").Int64()
				//KeyChainCarryInPrice, _ := battle.Get("KeyChainCarryInPrice").Int64()
				//CarryoutSelfPrice, _ := battle.Get("CarryoutSelfPrice").Int64()
				dtEventTime, _ := battle.Get("dtEventTime").String()
				ArmedForceId, _ := battle.Get("ArmedForceId").Int64()
				DurationS, _ := battle.Get("DurationS").Int64()
				KillCount, _ := battle.Get("KillCount").Int64()
				KillPlayerAICount, _ := battle.Get("KillPlayerAICount").Int64()
				KillAICount, _ := battle.Get("KillAICount").Int64()
				flowCalGainedPrice, _ := battle.Get("flowCalGainedPrice").Int64()
				battleS := fmt.Sprintf("[在%v使用%v在%v带出%v, 净收入%v, 击杀干员%v, 击杀人机干员%v, 击杀士兵%v, 生存时间%vmin, 撤离%v]\n",
					dtEventTime,
					armed.ArmedForceMap[int(ArmedForceId)],
					maps.MapIDMap[int(MapId)],
					FinalPrice,
					flowCalGainedPrice,
					KillCount,
					KillPlayerAICount,
					KillAICount,
					DurationS/60,
					int(EscapeFailReason),
				)
				m := message.Text(battleS)
				msgChan = append(msgChan, m)
			}
			if len(msgChan) == 1 {
				msgChan = append(msgChan, message.Text("没有战局"))
			}
			ctx.Send(msgChan)
		}

		if match, err := regexp.MatchString("^特勤处制造$", msg); err == nil && match {
			data, err := df.GetManufacture(fmt.Sprint(ctx.Event.UserID))
			if err != nil {
				ctx.SendChain(
					message.Reply(ctx.Event.MessageID),
					message.Text("失败: "+err.Error()),
				)
				return
			}
			placeData, _ := sonic.Get(data, "jData", "data", "data", "placeData")
			msgChain := []message.Segment{
				message.Reply(ctx.Event.MessageID),
			}
			placeDataLen := func() int { placeDataArr, _ := placeData.Array(); return len(placeDataArr) }()
			for i := 0; i < placeDataLen; i++ {
				Status, _ := placeData.Index(i).Get("Status").String()
				placeName, _ := placeData.Index(i).Get("placeName").String()
				if Status == "闲置中" {
					m := message.Text(fmt.Sprintf("%v: %v\n", placeName, Status))
					msgChain = append(msgChain, m)
				} else {
					leftTime, _ := placeData.Index(i).Get("leftTime").Int64()
					pushTime, _ := placeData.Index(i).Get("pushTime").Int64()
					objectId, _ := placeData.Index(i).Get("objectId").Int64()
					m := message.Text(fmt.Sprintf("%v: 生产%v中, 剩余时间%vmin, 制造完成时间%v\n",
						placeName,
						item.ItemMap[int(objectId)].ObjectName,
						leftTime/60,
						time.Unix(pushTime, 0).Format("2006-01-02 15:04:05"),
					))
					msgChain = append(msgChain, m)
				}
			}
			ctx.Send(msgChain)
		}

		//if match, err := regexp.MatchString("^最近带出$", msg); err == nil && match {
		//	pages := 100
		//	type dataS struct {
		//		data []byte
		//		ok   chan struct{}
		//	}
		//	data := make([]dataS, pages)
		//	for i := 0; i < pages; i++ {
		//		data[i] = dataS{
		//			ok: make(chan struct{}),
		//		}
		//		go func() {
		//			data[i].data, err = df.GetLoginPropsMoneyBattle(fmt.Sprint(ctx.Event.UserID), i+1, 2)
		//			data[i].ok <- struct{}{}
		//		}()
		//	}
		//	if err != nil {
		//		ctx.SendChain(
		//			message.Reply(ctx.Event.MessageID),
		//			message.Text("失败: "+err.Error()),
		//		)
		//		return
		//	}
		//
		//	msgChain := []message.Segment{
		//		message.Reply(ctx.Event.MessageID),
		//	}
		//	mtx := &sync.Mutex{}
		//	wg := &sync.WaitGroup{}
		//	for i := 0; i < pages; i++ {
		//		wg.Add(1)
		//		go func() {
		//			<-data[i].ok
		//			fmt.Println(i)
		//			itemArr, _ := sonic.Get(data[i].data, "jData", "data", "itemArr")
		//			fmt.Println(itemArr)
		//			itemArrLen := func() int {
		//				itemArrArr, _ := itemArr.Array()
		//				return len(itemArrArr)
		//			}()
		//			for j := 0; j < itemArrLen; j++ {
		//				iGoodsId, _ := itemArr.Index(j).Get("iGoodsId").Int64()
		//				AddOrReduce, _ := itemArr.Index(j).Get("AddOrReduce").String()
		//				Reason, _ := itemArr.Index(j).Get("Reason").String()
		//				Reason, _ = url.QueryUnescape(Reason)
		//				if item.ItemMap[int(iGoodsId)].Grade >= 6 && AddOrReduce == "+1" {
		//					dtEventTime, _ := itemArr.Index(j).Get("dtEventTime").String()
		//					m := message.Text(fmt.Sprintf("%v(%v)+1 at %v %v\n",
		//						item.ItemMap[int(iGoodsId)].ObjectName,
		//						item.ItemMap[int(iGoodsId)].AvgPrice,
		//						dtEventTime,
		//						Reason,
		//					))
		//					mtx.Lock()
		//					msgChain = append(msgChain, m)
		//					mtx.Unlock()
		//				}
		//			}
		//			wg.Done()
		//		}()
		//	}
		//	wg.Wait()
		//	ctx.Send(msgChain)
		//}

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
