package df

import (
	"strconv"
)

func GetBattle(qqid string, page int) ([]byte, error) {
	ck, err := GetCookie(qqid)
	if err != nil {
		return nil, err
	}
	qpage := (5*(page-1))/50 + 1
	return Reqest(ck, map[string]string{
		"iChartId":  "319386",
		"sIdeToken": "zMemOt",
		"type":      "4",
		"item":      "0,0,0,2201,0,0,0,75",
		"page":      strconv.Itoa(qpage),
	})
}
