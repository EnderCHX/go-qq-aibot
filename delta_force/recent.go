package df

func GetRecent(qqid string) ([]byte, error) {
	cookie, err := GetCookie(qqid)
	if err != nil {
		return nil, err
	}
	return Reqest(cookie, map[string]string{
		"iChartId":  "316969",
		"sIdeToken": "NoOapI",
		"method":    "dfm/center.recent.detail",
		"source":    "5",
		"param":     `{"resourceType":"sol"}`,
	})
}
