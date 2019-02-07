package shared

import (
	"github.com/byuoitav/common/nerr"
	sd "github.com/byuoitav/common/state/statedefinition"
	"github.com/byuoitav/common/v2/events"
)

//Cache is our state cache - it's meant to be a representation of the static indexes
type Cache interface {
	CheckAndStoreDevice(device sd.StaticDevice) (bool, sd.StaticDevice, *nerr.E)
	CheckAndStoreRoom(room sd.StaticRoom) (bool, sd.StaticRoom, *nerr.E)

	GetDeviceRecord(deviceID string) (sd.StaticDevice, *nerr.E)
	GetRoomRecord(roomID string) (sd.StaticRoom, *nerr.E)
	GetAllDeviceRecords() ([]sd.StaticDevice, *nerr.E)
	GetAllRoomRecords() ([]sd.StaticRoom, *nerr.E)

	StoreDeviceEvent(toSave sd.State) (bool, sd.StaticDevice, *nerr.E)
	StoreAndForwardEvent(event events.Event) (bool, *nerr.E)

	//GetDeviceManagerList() (int, []string, *nerr.E)
	GetCacheType() string
	GetCacheName() string
}
