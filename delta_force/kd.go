package df

func GetKd(qqid string) ([]byte, error) {
	ck, err := GetCookie(qqid)
	if err != nil {
		return nil, err
	}
	return Reqest(ck, map[string]string{
		"iChartId":  "317814",
		"sIdeToken": "QIRBwm",
	})
}
