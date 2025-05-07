package item

import (
	_ "embed"
	"github.com/bytedance/sonic"
)

//go:embed items.json
var itemsJsonBytes []byte

type Item struct {
	Id            int    `json:"id"`
	ObjectId      int    `json:"objectID"`
	ObjectName    string `json:"objectName"`
	Length        int    `json:"length"`
	Width         int    `json:"width"`
	Grade         int    `json:"grade"`
	Weight        string `json:"weight"`
	PrimaryClass  string `json:"primaryClass"`
	SecondClass   string `json:"secondClass"`
	SecondClassCN string `json:"secondClassCN"`
	Desc          string `json:"desc"`
	Pic           string `json:"pic"`
	PrePic        string `json:"prePic"`
	AvgPrice      int    `json:"avgPrice"`
	ProtectDetail struct {
		Durability   int `json:"durability"`
		ProtectLevel int `json:"protectLevel"`
		AimSpeed     struct {
			Percent      int    `json:"percent"`
			BatteryValue int    `json:"batteryValue"`
			BatteryColor string `json:"batteryColor"`
		} `json:"aimSpeed"`
		MoveSpeed struct {
			Percent      int    `json:"percent"`
			BatteryValue int    `json:"batteryValue"`
			BatteryColor string `json:"batteryColor"`
		} `json:"moveSpeed"`
		SoundEffect struct {
			Percent      int    `json:"percent"`
			BatteryValue int    `json:"batteryValue"`
			BatteryColor string `json:"batteryColor"`
		} `json:"soundEffect"`
		FaceMask struct {
			Value        string `json:"value"`
			BatteryValue int    `json:"batteryValue"`
			BatteryColor string `json:"batteryColor"`
		} `json:"faceMask"`
		RepairEfficiency string `json:"repairEfficiency"`
		ProtectArea      string `json:"protectArea"`
		DurableLoss      string `json:"durableLoss"`
	} `json:"protectDetail"`
	AccDetail struct {
		ControlSpeed  int `json:"controlSpeed"`
		QuickSeparate int `json:"quickSeparate"`
		Advantage     struct {
			Condition  string `json:"condition"`
			EffectList []struct {
				Value        string `json:"value"`
				BatteryValue int    `json:"batteryValue"`
				BatteryColor string `json:"batteryColor"`
			} `json:"effectList"`
		} `json:"advantage"`
		Disadvantage struct {
			Condition  string `json:"condition"`
			EffectList []struct {
				Value        string `json:"value"`
				BatteryValue int    `json:"batteryValue"`
				BatteryColor string `json:"batteryColor"`
			} `json:"effectList"`
		} `json:"disadvantage"`
	} `json:"accDetail"`
	PropsDetail struct {
		Type             string `json:"type"`
		PropsSource      string `json:"propsSource"`
		RepairPoints     int    `json:"repairPoints"`
		RepairArea       string `json:"repairArea"`
		RepairEfficiency string `json:"repairEfficiency"`
		ActiveTime       string `json:"activeTime"`
		ReplyEffect      string `json:"replyEffect"`
		HearEnhance      string `json:"hearEnhance"`
	} `json:"propsDetail"`
}

var ItemMap = make(map[int]Item)

func init() {
	var items []Item
	list, err := sonic.Get(itemsJsonBytes, "jData", "data", "data", "list")
	if err != nil {
		panic(err)
	}
	listRaw, err := list.MarshalJSON()
	if err != nil {
		panic(err)
	}
	sonic.Unmarshal(listRaw, &items)

	for _, item := range items {
		ItemMap[item.ObjectId] = item
	}
}
