package memorycache

import (
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	sd "github.com/byuoitav/common/state/statedefinition"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/shipwright/state/cache/shared"
	"github.com/robfig/cron"
)

type Memorycache struct {
	devicelock  sync.RWMutex
	deviceCache map[string]DeviceItemManager
	roomlock    sync.RWMutex
	roomCache   map[string]RoomItemManager

	cacheType string
	name      string

	pushCron *cron.Cron
}

func (c *Memorycache) GetCacheType() string {
	return c.cacheType
}

func (c *Memorycache) GetCacheName() string {
	return c.name
}

func (c *Memorycache) GetDeviceManagerList() (int, []string, *nerr.E) {
	toReturn := []string{}
	for k := range c.deviceCache {
		toReturn = append(toReturn, k)
	}

	return len(c.deviceCache), toReturn, nil
}

func (c *Memorycache) StoreAndForwardEvent(v events.Event) (bool, *nerr.E) {
	return shared.ForwardAndStoreEvent(v, c)
}

/*
	StoreDeviceEvent takes an event (key value) and stores the value in the field defined as key on a device.S
	Defer use to CheckAndStoreDevice for internal use, as there are significant speed gains.
*/
func (c *Memorycache) StoreDeviceEvent(toSave sd.State) (bool, sd.StaticDevice, *nerr.E) {
	if len(toSave.ID) < 1 {
		return false, sd.StaticDevice{}, nerr.Create("State must include device ID", "invaid-parameter")
	}

	c.devicelock.RLock()
	manager, ok := c.deviceCache[toSave.ID]
	c.devicelock.RUnlock()
	if !ok {
		log.L.Infof("Creating a new device manager for %v", toSave.ID)

		var err *nerr.E
		//we need to create a new manager and set it up
		manager, err = GetNewDeviceManager(toSave.ID)
		if err != nil {
			return false, sd.StaticDevice{}, err.Addf("couldn't store device event")
		}

		c.devicelock.Lock()
		c.deviceCache[toSave.ID] = manager
		c.devicelock.Unlock()
	}

	respChan := make(chan DeviceTransactionResponse, 1)

	//send a request to update
	manager.WriteRequests <- DeviceTransactionRequest{
		EventEdit:    true,
		Event:        toSave,
		ResponseChan: respChan,
	}

	//wait for a response
	resp := <-respChan

	if resp.Error != nil {
		return false, sd.StaticDevice{}, resp.Error.Addf("Couldn't store event %v.", toSave)
	}

	return resp.Changes, resp.NewDevice, nil
}

/*CheckAndStoreDevice takes a device, will check to see if there are deltas compared to the values in the map, and store any changes.

Bool returned denotes if there were any changes. True indicates that there were updates
*/
func (c *Memorycache) CheckAndStoreDevice(device sd.StaticDevice) (bool, sd.StaticDevice, *nerr.E) {
	if len(device.DeviceID) == 0 {
		return false, sd.StaticDevice{}, nerr.Create("Static Device must have an ID field to be loaded into the databaset", "invalid-device")
	}

	c.devicelock.RLock()
	manager, ok := c.deviceCache[device.DeviceID]
	c.devicelock.RUnlock()

	if !ok {
		var err *nerr.E
		manager, err = GetNewDeviceManager(device.DeviceID)
		if err != nil {
			return false, device, err.Addf("Couldn't check and store device")
		}

		c.devicelock.Lock()
		c.deviceCache[device.DeviceID] = manager
		c.devicelock.Unlock()
	}

	respChan := make(chan DeviceTransactionResponse, 1)

	//send a request to update
	manager.WriteRequests <- DeviceTransactionRequest{
		MergeDeviceEdit: true,
		MergeDevice:     device,
		ResponseChan:    respChan,
	}

	//wait for a response
	resp := <-respChan

	if resp.Error != nil {
		return false, sd.StaticDevice{}, resp.Error.Addf("Couldn't store device %v.", device)
	}

	shared.ForwardDevice(resp.NewDevice, resp.Changes, c)

	return resp.Changes, resp.NewDevice, nil
}

//GetDeviceRecord returns a device with the corresponding ID, if any is found in the memorycache
func (c *Memorycache) GetDeviceRecord(deviceID string) (sd.StaticDevice, *nerr.E) {

	manager, ok := c.deviceCache[deviceID]
	if !ok {
		return sd.StaticDevice{}, nil
	}

	respChan := make(chan sd.StaticDevice, 1)

	manager.ReadRequests <- respChan
	return <-respChan, nil
}

/*CheckAndStoreRoom takes a room, will check to see if there are deltas compared to the values in the map, and store any changes.

Bool returned denotes if there were any changes. True indicates that there were updates
Room returned contains ONLY the deltas.
*/
func (c *Memorycache) CheckAndStoreRoom(room sd.StaticRoom) (bool, sd.StaticRoom, *nerr.E) {
	if len(room.RoomID) == 0 {
		return false, sd.StaticRoom{}, nerr.Create("Static room must have a roomID to be compared and stored", "invalid-room")
	}

	manager, ok := c.roomCache[room.RoomID]
	if !ok {
		manager = GetNewRoomManager(room.RoomID)
	}

	respChan := make(chan RoomTransactionResponse, 1)

	//send a request to update
	manager.WriteRequests <- RoomTransactionRequest{
		MergeRoom:    room,
		ResponseChan: respChan,
	}

	//wait for a response
	resp := <-respChan

	if resp.Error != nil {
		return false, sd.StaticRoom{}, resp.Error.Addf("Couldn't store room %v.", room)
	}

	shared.ForwardRoom(resp.NewRoom, resp.Changes, c)

	return resp.Changes, resp.NewRoom, nil
}

func (c *Memorycache) PushAllDevices() {
	shared.PushAllDevices(c)
}

//GetRoomRecord returns a room
func (c *Memorycache) GetRoomRecord(roomID string) (sd.StaticRoom, *nerr.E) {
	manager, ok := c.roomCache[roomID]
	if !ok {
		return sd.StaticRoom{}, nil
	}

	respChan := make(chan sd.StaticRoom, 1)

	manager.ReadRequests <- respChan
	return <-respChan, nil
}

func (c *Memorycache) GetAllDeviceRecords() ([]sd.StaticDevice, *nerr.E) {
	toReturn := []sd.StaticDevice{}

	expected := len(c.deviceCache)
	ReadChannel := make(chan sd.StaticDevice, expected)

	c.devicelock.RLock()
	for _, v := range c.deviceCache {
		v.ReadRequests <- ReadChannel
	}
	c.devicelock.RUnlock()

	timeoutTimer := time.NewTimer(1 * time.Second)

	received := 0
	for {
		select {
		case <-timeoutTimer.C:
			log.L.Infof("ReadAll devices timed out..")
			return toReturn, nil
		case v := <-ReadChannel:
			toReturn = append(toReturn, v)
			received++
			if received >= expected {
				log.L.Debugf("Got all responses from the read all devices")
				return toReturn, nil
			}
		}
	}
}

func (c *Memorycache) GetAllRoomRecords() ([]sd.StaticRoom, *nerr.E) {
	toReturn := []sd.StaticRoom{}

	expected := len(c.deviceCache)
	ReadChannel := make(chan sd.StaticRoom, expected)

	c.roomlock.RLock()
	for _, v := range c.roomCache {
		v.ReadRequests <- ReadChannel
	}
	c.roomlock.RUnlock()

	timeoutTimer := time.NewTimer(1 * time.Second)

	received := 0
	for {
		select {
		case <-timeoutTimer.C:
			log.L.Infof("ReadAll rooms timed out..")
			return toReturn, nil
		case v := <-ReadChannel:
			toReturn = append(toReturn, v)
			received++
			if received >= expected {
				log.L.Debugf("Got all responses from the read all rooms")
				return toReturn, nil
			}
		}
	}
}

func (c *Memorycache) RemoveDevice(id string) *nerr.E {
	c.devicelock.Lock()
	manager, ok := c.deviceCache[id]
	if !ok {
		return nil
	}

	manager.KillChannel <- true

	delete(c.deviceCache, id)
	c.devicelock.Unlock()
	return nil
}

func (c *Memorycache) RemoveRoom(id string) *nerr.E {

	c.roomlock.Lock()
	manager, ok := c.roomCache[id]
	if !ok {
		return nil
	}

	manager.KillChannel <- true

	delete(c.roomCache, id)

	c.roomlock.Unlock()
	return nil
}

func (c *Memorycache) NukeRoom(id string) ([]string, *nerr.E) {
	er := c.RemoveRoom(id)
	if er != nil {
		return []string{}, er.Addf("Couldn't nuke room")
	}

	toDelete := []string{}
	c.devicelock.RLock()
	for k := range c.deviceCache {
		if strings.HasPrefix(k, id) {
			toDelete = append(toDelete, k)
		}
	}
	c.devicelock.RUnlock()

	for i := range toDelete {
		er = c.RemoveDevice(toDelete[i])
		if er != nil {
			return []string{}, er.Addf("Couldn't nuke room")
		}
	}

	return toDelete, nil
}
