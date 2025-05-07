package df

func GetAchievements(qqid string) ([]byte, error) {
	ck, err := GetCookie(qqid)
	if err != nil {
		return nil, err
	}
	return Reqest(ck, map[string]string{
		"iChartId":  "384918",
		"sIdeToken": "mbq5GZ",
		"method":    "dist.contents",
		"source":    "5",
		"param":     `{"distType":"bannerManage","contentType":"secretDay"}`,
	})
}
