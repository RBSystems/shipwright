package store

import (
	"fmt"
	"time"

	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
)

type alertStore struct {
	inChannel      chan structs.Alert
	requestChannel chan alertRequest

	store map[string]structs.Alert
	//	configuration Config
}

//AlertRequest is submitted to the store to retrieve an alert from it.
type alertRequest struct {
	AlertID      string
	ResponseChan chan alertResponse
	All          bool
}

//AlertRepsonse should always have the error checked before retrieving the alert
type alertResponse struct {
	Error *nerr.E
	Alert []structs.Alert
}

var ZeroTime = time.Time{}

var store alertStore

func init() {
	store := alertStore{
		inChannel:      make(chan structs.Alert, 1000),
		requestChannel: make(chan alertRequest, 1000),
		store:          map[string]structs.Alert{},
	}

	go store.run()
}

//PutAlert adds an alert to the store.
//Do we want to wait for confirmation?
func (a *alertStore) putAlert(alert structs.Alert) (string, *nerr.E) {

	//check to make sure we have a time
	if alert.AlertStartTime.IsZero() {
		alert.AlertStartTime = time.Now()
	}

	//Check to make sure we have an ID
	if alert.AlertID == "" {
		//we need to generate
		alert.AlertID = GenerateID(alert)
	}

	a.inChannel <- alert

	return alert.AlertID, nil
}

func (a *alertStore) getAlert(id string) (structs.Alert, *nerr.E) {

	//make our request
	respChan := make(chan alertResponse, 1)

	a.requestChannel <- alertRequest{
		AlertID:      id,
		ResponseChan: respChan,
	}

	resp := <-respChan

	if len(resp.Alert) > 0 {
		return resp.Alert[0], resp.Error
	}

	return structs.Alert{}, resp.Error
}

//NOT SAFE FOR CONCURRENT ACCESS. DO NOT USE OUTSIDE OF run()
func (a *alertStore) resolveAlert(alertID string, resInfo structs.ResolutionInfo) *nerr.E {
	//we remove it from the store, and ship it off to the persistance stuff.
	//we should check to see if it already exists
	if v, ok := a.store[alertID]; ok {

		//it's there, lets get it, mark it as resolved.
		v.Resolved = true
		v.ResolutionInfo = resInfo
		v.AlertID = v.AlertID + v.AlertStartTime.Format(time.RFC3339) //change the ID so it's unique

		delete(a.store, alertID)

		//submit for persistance in the resolved alert

		//remove from persistance in the active alerts array

	} else {
		return nerr.Create("Unkown alert "+alertID, "not-found")
	}

	return nil
}

func (a *alertStore) run() {

	for {
		select {
		case al := <-a.inChannel:
			a.storeAlert(al)
		case req := <-a.requestChannel:
			a.handleRequest(req)
		}
	}
}

//NOT SAFE FOR CONCURRENT ACCESS. DO NOT USE OUTSIDE OF run()
func (a *alertStore) storeAlert(alert structs.Alert) {

	//we should check to see if it already exists
	if v, ok := a.store[alert.AlertID]; ok {

		if len(alert.Message) > 0 &&
			(len(v.MessageLog) == 0 || v.MessageLog[len(v.MessageLog)-1] != alert.Message) {

			v.MessageLog = append(v.MessageLog, alert.Message)
			v.Message = alert.Message
		}

		v.Active = alert.Active
		v.AlertLastUpdateTime = time.Now()

		a.store[v.AlertID] = v //store it back in

	} else {

		//we store it.
		alert.AlertLastUpdateTime = time.Now()
		if len(alert.Message) > 0 {
			alert.MessageLog = append(v.MessageLog, alert.Message)
		}

		a.store[alert.AlertID] = alert
	}

	//submit to be updated in persistance.

	//check rules for actions
}

//NOT SAFE FOR CONCURRENT ACCESS. DO NOT USE OUTSIDE OF run()
func (a *alertStore) handleRequest(req alertRequest) {
	toReturn := []structs.Alert{}

	if req.All {
		for _, v := range a.store {
			toReturn = append(toReturn, v)
		}

	} else {
		if v, ok := a.store[req.AlertID]; ok {

			toReturn = append(toReturn, v)
			req.ResponseChan <- alertResponse{
				Error: nil,
				Alert: toReturn,
			}
		} else {
			req.ResponseChan <- alertResponse{
				Error: nerr.Create(fmt.Sprintf("No Alert for id %v", req.AlertID), "not-found"),
				Alert: toReturn,
			}
		}

	}

	req.ResponseChan <- alertResponse{
		Error: nil,
		Alert: toReturn,
	}
}
