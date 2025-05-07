package df

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

type ReqBody struct {
	IChartId    int    `json:"iChartId"`
	ISubChartId int    `json:"iSubChartId"`
	SIdeToken   string `json:"sIdeToken"`
	E_code      int    `json:"e_code"`
	G_code      int    `json:"g_code"`
	Eas_url     string `json:"eas_url"`
	Eas_refer   string `json:"eas_refer"`
	SMiloTag    string `json:"sMiloTag"`
	Method      string `json:"method"`
	Source      int    `json:"source"`
	Param       string `json:"param"`
}

func (r *ReqBody) Read(p []byte) (n int, err error) {
	//p, err = json.Marshal(r)
	//n = len(p)
	return
}

func Reqest(cookie string, param map[string]string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Cookie", cookie).
		SetHeader("accept-language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2").
		SetQueryParams(param).
		Post("https://comm.ams.game.qq.com/ide/")
	sMsg := "succ"
	if data, err := sonic.Get(resp.Body(), "sMsg"); err != nil || func() bool {
		sMsg, _ = data.String()
		if sMsg == "succ" || sMsg == "ok" {
			return false
		} else {
			return true
		}
	}() {
		if err != nil {
			return resp.Body(), err
		} else {
			return resp.Body(), fmt.Errorf("error: %s", sMsg)
		}
	}
	return resp.Body(), err
}
