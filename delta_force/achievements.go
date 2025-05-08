package df

func GetAchievements(qqid string) ([]byte, error) {
	ck, err := GetCookie(qqid)
	if err != nil {
		return nil, err
	}
	return Reqest(ck, map[string]string{
		"iChartId":  "316969",
		"sIdeToken": "NoOapI",
		"method":    "dfm/center.person.resource",
		"source":    "5",
		"param":     `{"resourceType":"sol","seasonid":[1,2,3,4],"isAllSeason":true}`,
	})
}
