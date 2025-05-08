package df

func GetManufacture(qqid string) ([]byte, error) {
	ck, err := GetCookie(qqid)
	if err != nil {
		return nil, err
	}
	return Reqest(ck, map[string]string{
		"iChartId":  "365589",
		"sIdeToken": "bQaMCQ",
		"source":    "5",
	})
}
