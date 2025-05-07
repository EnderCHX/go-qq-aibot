package df

func GetPassword(qqid string) ([]byte, error) {
	cookie, err := GetCookie(qqid)
	if err != nil {
		return nil, err
	}
	return Reqest(cookie, map[string]string{
		"iChartId":  "384918",
		"sIdeToken": "mbq5GZ",
		"method":    "dist.contents",
		"source":    "5",
		"param":     `{"distType":"bannerManage","contentType":"secretDay"}`,
	})
}
