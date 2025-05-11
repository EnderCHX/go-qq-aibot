package df

import (
	"strconv"
)

// 查询类型（1：登录，2：道具，3：货币,4:战绩）
func GetLoginPropsMoneyBattle(qqid string, page int, type_ int) ([]byte, error) {
	ck, err := GetCookie(qqid)
	if err != nil {
		return nil, err
	}
	qpage := (5*(page-1))/50 + 1
	return Reqest(ck, map[string]string{
		"iChartId":  "319386",
		"sIdeToken": "zMemOt",
		"type":      strconv.Itoa(type_),
		"page":      strconv.Itoa(qpage),
	})
}
