package maps

import _ "embed"

//go:embed map.json
var mapJsonBytes []byte

var MapIDMap = map[int]string{
	2201: "零号大坝-常规",
	2202: "零号大坝-机密",
	1901: "长工溪谷-常规",
	3901: "航天基地-机密",
	3902: "航天基地-绝密",
	8102: "巴克什-机密",
}
