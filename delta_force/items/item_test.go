package item

import (
	"fmt"
	"github.com/bytedance/sonic"
	"testing"
)

func TestItem(t *testing.T) {
	var items []Item
	list, err := sonic.Get(itemsJsonBytes, "jData", "data", "data", "list")
	if err != nil {
		t.Error(err)
	}
	listRaw, err := list.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	sonic.Unmarshal(listRaw, &items)
	//fmt.Println(items)
	itemsMap := map[int]Item{}
	for _, item := range items {
		itemsMap[item.ObjectId] = item
	}
	fmt.Println(itemsMap)
}
