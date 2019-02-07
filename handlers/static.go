package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/shipwright/helpers"
	"github.com/labstack/echo"
)

// GetAllStaticDeviceRecords returns a list of all the static device records
func GetAllStaticDeviceRecords(context echo.Context) error {
	log.L.Debugf("%s Starting GetAllStaticDeviceRecords...", helpers.StaticTag)

	sDevices, err := helpers.GetAllStaticDeviceRecords()
	if err != nil {
		msg := fmt.Sprintf("failed to get all static device records : %s", err.Error())
		log.L.Errorf("%s %s", helpers.StaticTag, msg)
		return context.JSON(http.StatusInternalServerError, err)
	}

	log.L.Debugf("%s Successfully got all static device records!", helpers.StaticTag)
	return context.JSON(http.StatusOK, sDevices)
}

// GetAllStaticRoomRecords returns a list of all the static room records
func GetAllStaticRoomRecords(context echo.Context) error {
	log.L.Debugf("%s Starting GetAllStaticRoomRecords...", helpers.StaticTag)

	sRooms, err := helpers.GetAllStaticRoomRecords()
	if err != nil {
		msg := fmt.Sprintf("failed to get all static room records : %s", err.Error())
		log.L.Errorf("%s %s", helpers.StaticTag, msg)
		return context.JSON(http.StatusInternalServerError, err)
	}

	log.L.Debugf("%s Successfully got all static room records!", helpers.StaticTag)
	return context.JSON(http.StatusOK, sRooms)
}