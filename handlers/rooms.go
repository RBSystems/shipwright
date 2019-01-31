package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/shipwright/helpers"
	"github.com/labstack/echo"
)

// AddRoom adds a room to the database
func AddRoom(context echo.Context) error {
	log.L.Debugf("%s Starting AddRoom...", helpers.RoomsTag)

	// get information from the context
	roomID := context.Param("room")
	username := getUsernameString(context)

	var room structs.Room
	err := context.Bind(&room)
	if err != nil {
		msg := fmt.Sprintf("failed to bind request body for %s : %s", roomID, err.Error())
		log.L.Errorf("%s %s", helpers.RoomsTag, msg)
		return context.JSON(http.StatusBadRequest, err)
	}

	// call helper function
	result, ne := helpers.AddRoom(roomID, room)
	if ne != nil {
		log.L.Errorf("%s %s", helpers.RoomsTag, ne.Error())
		return context.JSON(http.StatusInternalServerError, result)
	}

	helpers.CreateAndAddRoomChange(room, helpers.AddAction, username)

	log.L.Debugf("%s The room %s was successfully created!", helpers.RoomsTag, roomID)
	return context.JSON(http.StatusOK, result)
}

// AddMultipleRooms adds a set of rooms to the database
func AddMultipleRooms(context echo.Context) error {
	log.L.Debugf("%s Starting AddMultipleRooms...", helpers.RoomsTag)

	// get information from the context
	username := getUsernameString(context)

	var rooms []structs.Room

	err := context.Bind(&rooms)
	if err != nil {
		msg := fmt.Sprintf("failed to bind request body for multiple rooms : %s", err.Error())
		log.L.Errorf("%s %s", helpers.RoomsTag, msg)
		return context.JSON(http.StatusBadRequest, err)
	}

	var results []helpers.DBResponse
	// call helper function as we iterate
	for _, r := range rooms {
		res, ne := helpers.AddRoom(r.ID, r)
		if ne != nil {
			log.L.Errorf("%s %s", helpers.RoomsTag, ne.Error())
		}

		helpers.CreateAndAddRoomChange(r, helpers.AddAction, username)
		results = append(results, res)
	}

	log.L.Debugf("%s The rooms were successfully created!", helpers.RoomsTag)
	return context.JSON(http.StatusOK, results)
}

// GetRoom gets a room from the database based on the given ID
func GetRoom(context echo.Context) error {
	log.L.Debugf("%s Starting GetRoom...", helpers.RoomsTag)

	// get information from the context
	roomID := context.Param("room")

	room, err := helpers.GetRoom(roomID)
	if err != nil {
		msg := fmt.Sprintf("failed to get the room %s : %s", roomID, err.Error())
		log.L.Errorf("%s %s", helpers.RoomsTag, msg)
		return context.JSON(http.StatusInternalServerError, err)
	}

	log.L.Debugf("%s Successfully found the room %s!", helpers.RoomsTag, roomID)
	return context.JSON(http.StatusOK, room)
}

// GetAllRooms gets all rooms from the database
func GetAllRooms(context echo.Context) error {
	log.L.Debugf("%s Starting GetAllRooms...", helpers.RoomsTag)

	rooms, err := helpers.GetAllRooms()
	if err != nil {
		msg := fmt.Sprintf("failed to get all rooms : %s", err.Error())
		log.L.Errorf("%s %s", helpers.RoomsTag, msg)
		return context.JSON(http.StatusBadRequest, err)
	}

	log.L.Debugf("%s Successfully got all rooms!", helpers.RoomsTag)
	return context.JSON(http.StatusOK, rooms)
}

// GetRoomsByBuilding gets all rooms in a single building from the database, based on a given building ID
func GetRoomsByBuilding(context echo.Context) error {
	log.L.Debugf("%s Starting GetRoomsByBuilding...", helpers.RoomsTag)

	// get information from the context
	buildingID := context.Param("building")

	rooms, err := helpers.GetRoomsByBuilding(buildingID)
	if err != nil {
		msg := fmt.Sprintf("failed to get all of the rooms in the building %s : %s", buildingID, err.Error())
		log.L.Errorf("%s %s", helpers.RoomsTag, msg)
		return context.JSON(http.StatusInternalServerError, err)
	}

	log.L.Debugf("%s Successfully got all rooms in the building %s", helpers.RoomsTag, buildingID)
	return context.JSON(http.StatusOK, rooms)
}

// UpdateRoom updates a room in the database
func UpdateRoom(context echo.Context) error {
	log.L.Debugf("%s Starting UpdateRoom...", helpers.RoomsTag)

	// get information from the context
	roomID := context.Param("room")
	username := getUsernameString(context)

	var room structs.Room
	err := context.Bind(&room)
	if err != nil {
		msg := fmt.Sprintf("failed to bind request body for %s : %s", roomID, err.Error())
		log.L.Errorf("%s %s", helpers.RoomsTag, msg)
		return context.JSON(http.StatusBadRequest, err)
	}

	// call the helper function
	result, ne := helpers.UpdateRoom(roomID, room)
	if ne != nil {
		log.L.Errorf("%s %s", helpers.RoomsTag, ne.Error())
		return context.JSON(http.StatusInternalServerError, result)
	}

	helpers.CreateAndAddRoomChange(room, helpers.UpdateAction, username)

	log.L.Debugf("%s The room %s was successfully updated!", helpers.RoomsTag, roomID)
	return context.JSON(http.StatusOK, result)
}

// UpdateMultipleRooms updates a set of rooms in the database
func UpdateMultipleRooms(context echo.Context) error {
	log.L.Debugf("%s Starting UpdateMultipleRooms...", helpers.RoomsTag)

	// get information from the context
	username := getUsernameString(context)

	var rooms map[string]structs.Room

	err := context.Bind(&rooms)
	if err != nil {
		msg := fmt.Sprintf("failed to bind request body for multiple buildings : %s", err.Error())
		log.L.Errorf("%s %s", helpers.RoomsTag, msg)
		return context.JSON(http.StatusBadRequest, err)
	}

	var results []helpers.DBResponse
	// call helper function as we iterate
	for id, room := range rooms {
		res, ne := helpers.UpdateRoom(id, room)
		if ne != nil {
			log.L.Errorf("%s %s", helpers.BuildingsTag, ne.Error())
		}

		helpers.CreateAndAddRoomChange(room, helpers.UpdateAction, username)
		results = append(results, res)
	}

	log.L.Debugf("%s The rooms were successfully updated!", helpers.RoomsTag)
	return context.JSON(http.StatusOK, results)
}

// DeleteRoom deletes a room from the database based on the given ID
func DeleteRoom(context echo.Context) error {
	log.L.Debugf("%s Starting DeleteRoom...", helpers.RoomsTag)

	// get information from the context
	roomID := context.Param("room")
	username := getUsernameString(context)

	room := structs.Room{ID: roomID}

	// call helper function
	result, ne := helpers.DeleteRoom(roomID)
	if ne != nil {
		log.L.Errorf("%s %s", helpers.RoomsTag, ne.Error())
		return context.JSON(http.StatusInternalServerError, result)
	}

	helpers.CreateAndAddRoomChange(room, helpers.DeleteAction, username)

	log.L.Debugf("%s The room %s was successfully deleted!", helpers.RoomsTag, roomID)
	return context.JSON(http.StatusOK, result)
}

// GetRoomConfigurations returns a list of possible room configurations
func GetRoomConfigurations(context echo.Context) error {
	log.L.Debugf("%s Starting GetRoomConfigurations...", helpers.RoomsTag)

	configurations, err := helpers.GetRoomConfigurations()
	if err != nil {
		msg := fmt.Sprintf("failed to get all room configurations : %s", err.Error())
		log.L.Errorf("%s %s", helpers.RoomsTag, msg)
		return context.JSON(http.StatusInternalServerError, err)
	}

	log.L.Debugf("%s Successfully got all room configurations!", helpers.RoomsTag)
	return context.JSON(http.StatusOK, configurations)
}

// GetRoomDesignations returns a list of possible room designations
func GetRoomDesignations(context echo.Context) error {
	log.L.Debugf("%s Starting GetRoomDesignations...", helpers.RoomsTag)

	designations, err := helpers.GetRoomDesignations()
	if err != nil {
		msg := fmt.Sprintf("failed to get all room designations : %s", err.Error())
		log.L.Errorf("%s %s", helpers.RoomsTag, msg)
		return context.JSON(http.StatusInternalServerError, err)
	}

	log.L.Debugf("%s Successfully got all room designations!", helpers.RoomsTag)
	return context.JSON(http.StatusOK, designations)
}