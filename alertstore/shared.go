package alertstore

import (
	"fmt"

	"github.com/byuoitav/common/structs"
)

const (
	Init      = "init"
	Interface = "interface"
	Cache     = "cache"
)

//GenerateID will give you the generated id back
func GenerateID(a structs.Alert) string {
	return fmt.Sprintf("%v^%v^%v", a.DeviceID, a.Type, a.Category)
}
