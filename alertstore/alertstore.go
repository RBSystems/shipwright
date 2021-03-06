package alertstore

import (
	"context"
	"fmt"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/shipwright/actions"
	"github.com/byuoitav/shipwright/actions/actionctx"
	"github.com/byuoitav/shipwright/alertstore/alertcache"
	"github.com/byuoitav/shipwright/alertstore/persist"
	"github.com/byuoitav/shipwright/socket"
)

type alertStore struct {
	inChannel         chan structs.Alert
	issueEditChannel  chan structs.RoomIssue
	requestChannel    chan issueRequest
	resolutionChannel chan resolutionRequest

	actionManager *actions.ActionManager
}

//issueRequest is submitted to the store to retrieve an alert from it. You can get alerts by ID, RoomIssue or Room
type issueRequest struct {
	RoomIssueID  string
	ResponseChan chan issueResponse
	All          bool
}

type resolutionRequest struct {
	RoomIssue      string //resolve a room issue
	ResolutionInfo structs.ResolutionInfo

	Partial  bool     //resolve a room issue
	AlertIDs []string //for use if partial
}

//issueResponse should always have the error checked before retrieving the alert
type issueResponse struct {
	Error      *nerr.E
	RoomIssues []structs.RoomIssue
	Alert      structs.Alert
}

//ZeroTime .
var ZeroTime = time.Time{}

var store *alertStore

//InitializeAlertStore sets up the alertstore
func InitializeAlertStore(a *actions.ActionManager) {
	store = &alertStore{
		inChannel:         make(chan structs.Alert, 1000),
		issueEditChannel:  make(chan structs.RoomIssue, 1000),
		requestChannel:    make(chan issueRequest, 1000),
		resolutionChannel: make(chan resolutionRequest, 1000),
		actionManager:     a,
	}

	go store.run()

	issues, err := persist.GetAllActiveIssuesFromPersist()
	if err != nil {
		log.L.Errorf("Couldn't get all active alerts: %v", err.Error())
	}

	log.L.Infof("Initializing alert store with %v alerts", len(issues))
	for i := range issues {

		for _, v := range issues[i].Alerts {
			v.Source = Init
			v, err = AddRoomInformationToAlert(v)

			if err != nil {
				log.L.Warnf("Problem adding room info to alert %v", err.Error())
			}
			log.L.Debugf("Storing alert %v", v)
			store.inChannel <- v
		}

		store.issueEditChannel <- issues[i]
	}

	log.L.Infof("Alert store initialized with %v issues", len(issues))
}

// Status .
type Status struct {
	QueueStatus map[string]ChannelStatus
}

// ChannelStatus is the status of the channel
type ChannelStatus struct {
	Cap int `json:"cap"`
	Len int `json:"len"`
}

func (a *alertStore) getQueueStatus() Status {
	return Status{
		QueueStatus: map[string]ChannelStatus{
			"in-queue": ChannelStatus{
				Cap: cap(a.inChannel),
				Len: len(a.inChannel),
			},
			"request-queue": ChannelStatus{
				Cap: cap(a.requestChannel),
				Len: len(a.requestChannel),
			},
			"issue-queue": ChannelStatus{
				Cap: cap(a.issueEditChannel),
				Len: len(a.issueEditChannel),
			},
			"resolution-queue": ChannelStatus{
				Cap: cap(a.resolutionChannel),
				Len: len(a.resolutionChannel),
			},
			// "room-aggregation-queue": ChannelStatus{
			// 	Cap: cap(roomAggsInChannel),
			// 	Len: len(roomAggsInChannel),
			// },
		},
	}
}

func (a *alertStore) setRoomIssueInfo(issue structs.RoomIssue) *nerr.E {

	//we must have a room issue ID
	if issue.RoomIssueID == "" {
		return nerr.Create("unable to set the room issue info, must have an issue ID", "insufficient-info")
	}
	if issue.Resolved {
		return nerr.Create("Can't respolve using setRoomIssueInfo, use ResolveIssue", "invalid-args")
	}

	//otherwise we ship it in to be processed
	log.L.Infof("Editing the room issue info for %v", issue.RoomIssueID)
	log.L.Infof("IncidentIDS: %v", issue.IncidentID)

	a.issueEditChannel <- issue
	return nil
}

//putAlert adds an alert to the store.
func (a *alertStore) putAlert(alert structs.Alert) (string, *nerr.E) {
	//check to make sure we have a time
	if alert.AlertStartTime.IsZero() {
		alert.AlertStartTime = time.Now()
	}

	//Check to make sure we have an ID
	if alert.AlertID == "" {
		//we need to generate
		alert.AlertID = GenerateAlertID(alert)
	}

	log.L.Debugf("Adding alert %v for device %v", alert.AlertID, alert.DeviceID)

	a.inChannel <- alert
	return alert.AlertID, nil
}

func (a *alertStore) resolveRoomIssue(resolutionInfo structs.ResolutionInfo, roomIssue string) *nerr.E {
	a.resolutionChannel <- resolutionRequest{
		RoomIssue:      roomIssue,
		ResolutionInfo: resolutionInfo,
	}

	return nil
}

func (a *alertStore) resolvePartialRoomIssue(resolutionInfo structs.ResolutionInfo, roomIssue string, alertIDs []string) *nerr.E {

	a.resolutionChannel <- resolutionRequest{
		RoomIssue:      roomIssue,
		ResolutionInfo: resolutionInfo,
		Partial:        true,
		AlertIDs:       alertIDs,
	}

	return nil
}

func (a *alertStore) getAllIssues() ([]structs.RoomIssue, *nerr.E) {
	log.L.Infof("Getting all room issues")

	//make our request
	respChan := make(chan issueResponse, 1)

	a.requestChannel <- issueRequest{
		All:          true,
		ResponseChan: respChan,
	}

	resp := <-respChan

	return resp.RoomIssues, resp.Error

}
func (a *alertStore) getRoomIssue(id string) (structs.RoomIssue, *nerr.E) {
	log.L.Infof("Getting room issue %v", id)

	//make our request
	respChan := make(chan issueResponse, 1)

	a.requestChannel <- issueRequest{
		RoomIssueID:  id,
		ResponseChan: respChan,
	}

	resp := <-respChan

	if len(resp.RoomIssues) < 1 {
		return structs.RoomIssue{}, resp.Error
	}

	return resp.RoomIssues[0], resp.Error

}

func (a *alertStore) run() {
	log.L.Infof("running alert store")

	for {
		log.L.Debugf("Waiting for event")
		select {
		case al := <-a.inChannel:
			log.L.Debugf("Storing alert")
			a.storeAlert(al)
		case is := <-a.issueEditChannel:
			log.L.Debugf("edit issue information")
			a.editIssueInformation(is)
		case req := <-a.requestChannel:
			a.handleRequest(req)
		case req := <-a.resolutionChannel:
			a.resolveIssue(req.ResolutionInfo, req.RoomIssue, req.Partial, req.AlertIDs)
		}
	}
}

//NOT SAFE FOR CONCURRENT ACCESS. DO NOT USE OUTSIDE OF run()
func (a *alertStore) editIssueInformation(issue structs.RoomIssue) *nerr.E {
	//check to see if it exists
	v, err := alertcache.GetAlertCache("default").GetIssue(issue.RoomIssueID)
	if err != nil {
		if err.Type == alertcache.NotFound {
			log.L.Errorf("Trying to edit room issue that doesn't exist: %v", issue.RoomIssueID)
			return nerr.Create("Trying to edit room issue that doesn't exist", "bad-id")
		}

		log.L.Errorf("Couldn't get room issue %v from cache %v", issue.RoomIssueID, err.Error())
		return err
	}

	//v is our deal, we need to combine
	i, changes := combineIssues(issue, v)
	if changes {
		persist.GetElkAlertPersist().StoreIssue(i, false, false)
		alertcache.GetAlertCache("default").PutIssue(i)

		a.runIssueActions(i)
		socket.GetManager().WriteToSockets(i)
	}

	return nil
}

//NOT SAFE FOR CONCURRENT ACCESS. DO NOT USE OUTSIDE OF run()
func (a *alertStore) resolveIssue(resInfo structs.ResolutionInfo, roomIssue string, partial bool, alertIDs []string) *nerr.E {
	log.L.Infof("Resolving issue %v", roomIssue)
	v, err := alertcache.GetAlertCache("default").GetIssue(roomIssue)
	if err == nil {
		if partial {
			log.L.Infof("Resolving partial issue %v. Resolving alerts: %v", roomIssue, alertIDs)

			//we need to copy the RoomIssue
			newRoomIssue := v

			//copy some of our slices
			copy(newRoomIssue.RoomTags, v.RoomTags)
			// copy(newRoomIssue.AlertTypes, v.AlertTypes)
			// copy(newRoomIssue.AlertCategories, v.AlertCategories)
			// copy(newRoomIssue.ActiveAlertTypes, v.ActiveAlertTypes)
			// copy(newRoomIssue.ActiveAlertCategories, v.ActiveAlertCategories)
			newRoomIssue.Alerts = []structs.Alert{}
			copy(newRoomIssue.RoomIssueResponses, v.RoomIssueResponses)
			copy(newRoomIssue.NotesLog, v.NotesLog)

			keepAlerts := []structs.Alert{}

			//go through and parse out the alerts from old and move them to the new one.
			for i := range v.Alerts {
				found := false
				for j := range alertIDs {
					if alertIDs[j] == v.Alerts[i].AlertID {
						//they match, add it to our new boy
						newRoomIssue.Alerts = append(newRoomIssue.Alerts, v.Alerts[i])
						found = true
					}
				}
				if !found {
					//add it to old
					keepAlerts = append(keepAlerts, v.Alerts[i])
				}
			}

			v.Alerts = keepAlerts

			newRoomIssue.ResolutionInfo = resInfo
			newRoomIssue.Resolved = true

			v.CalculateAggregateInfo()
			newRoomIssue.CalculateAggregateInfo()

			log.L.Debugf("partial resolve - new issue: %+v, old issue: %+v", newRoomIssue, v)

			//update in cache
			alertcache.GetAlertCache("default").PutIssue(v)

			//persist to elk for storage
			persist.GetElkAlertPersist().StoreIssue(v, true, false)
			socket.GetManager().WriteToSockets(v)
			a.runIssueActions(v)

			persist.GetElkAlertPersist().StoreIssue(newRoomIssue, true, true)
			socket.GetManager().WriteToSockets(newRoomIssue)
			a.runIssueActions(newRoomIssue)
		} else {
			log.L.Infof("Resolving full issue %v", roomIssue)
			err := alertcache.GetAlertCache("default").DeleteIssue(roomIssue)
			if err != nil {
				return err.Addf("couldn't resolve issue: %v", err.Error())
			}

			//it's there, lets get it, mark it as resolved.
			v.Resolved = true
			v.ResolutionInfo = resInfo

			// close each of the alerts
			for i := range v.Alerts {
				v.Alerts[i].Active = false
				a.runAlertActions(v.Alerts[i])
			}

			//submit for persistence
			v.CalculateAggregateInfo()

			persist.GetElkAlertPersist().StoreIssue(v, true, true)
			a.runIssueActions(v)
			socket.GetManager().WriteToSockets(v)
		}
	} else if err.Type == alertcache.NotFound {
		log.L.Errorf("%v", nerr.Create("Unkown room issue "+roomIssue, "not-found"))
	} else {
		log.L.Errorf("%v", err.Addf("couldn't resolve room issue"))
	}

	return nil
}

//NOT SAFE FOR CONCURRENT ACCESS. DO NOT USE OUTSIDE OF run()
func (a *alertStore) storeAlert(alert structs.Alert) {
	log.L.Debugf("Storing alert %v", alert.AlertID)

	//get the issue associated with the alert
	issueID := GetIssueIDFromAlertID(alert.AlertID)

	//Do we need to change the roomaggregate info?
	//roomAggregateChange := false

	//we should check to see if the room already has an issue associated with it.
	issue, err := alertcache.GetAlertCache("default").GetIssue(issueID)
	if err == nil {

		//we need to check to see if this alert exists on the issuecheck to see if this alert exists on the issue
		var v structs.Alert
		var ok bool

		var indx int
		for i := range issue.Alerts {
			if issue.Alerts[i].AlertID == alert.AlertID {
				ok = true
				v = issue.Alerts[i]
				indx = i

				break
			}
		}

		if ok {
			//check to see if our last update time is non-blank and before the one already in the cache, if so we don't do anything
			if !alert.AlertLastUpdateTime.IsZero() && alert.AlertLastUpdateTime.Before(v.AlertLastUpdateTime) {
				log.L.Infof("Alert: %v is out of date from cache.", alert.AlertID)
				//check if it's an init
				if alert.Source == Init {
					//create run the actions based on the alert in storage - since that's more up to date
					a.runAlertInitActions(v)
				}
				return
			}

			if len(alert.Message) > 0 &&
				(len(v.MessageLog) == 0 || v.MessageLog[len(v.MessageLog)-1] != alert.Message) {

				v.MessageLog = append(v.MessageLog, alert.Message)
				v.Message = alert.Message
			}

			var runInit bool
			if !alert.Active && v.Active {
				if alert.AlertEndTime.IsZero() {
					v.AlertEndTime = time.Now()
				} else {
					v.AlertEndTime = alert.AlertEndTime
				}
			} else if alert.Active && !v.Active {
				runInit = true
			}

			v.Active = alert.Active

			v.AlertLastUpdateTime = time.Now()

			if alert.SystemType != "" {
				v.SystemType = alert.SystemType
			}

			if len(alert.DeviceTags) >= 0 {
				v.DeviceTags = alert.DeviceTags
			}
			if len(alert.RoomTags) >= 0 {
				v.RoomTags = alert.RoomTags
			}
			if len(alert.AlertTags) >= 0 {
				v.AlertTags = alert.AlertTags
			}

			if alert.Source == Init || runInit {
				//create run the actions based on the alert in storage - since that's more up to date
				a.runAlertInitActions(v)
			}

			issue.Alerts[indx] = v

			alert = v
		} else {
			//roomAggregateChange = true
			//we store it.
			alert.AlertLastUpdateTime = time.Now()

			//set the start time
			if alert.AlertStartTime.IsZero() {
				alert.AlertStartTime = time.Now()
			}

			if len(alert.Message) > 0 {
				alert.MessageLog = append(v.MessageLog, alert.Message)
			}

			if alert.Active {
				//run the iniitialization actions thing
				a.runAlertInitActions(alert)
			}
			issue.Alerts = append(issue.Alerts, alert)
		}

	} else if err.Type == alertcache.NotFound {
		//issue didn't exist at all.
		//roomAggregateChange = true

		//generate the new roomIssue.
		alert.AlertLastUpdateTime = time.Now()

		//set the start time
		if alert.AlertStartTime.IsZero() {
			alert.AlertStartTime = time.Now()
		}

		if len(alert.Message) > 0 {
			alert.MessageLog = append(alert.MessageLog, alert.Message)
		}

		if alert.Active {
			//run the iniitialization actions thing
			a.runAlertInitActions(alert)
		}

		issue = structs.RoomIssue{
			RoomIssueID:     GetIssueIDFromAlertID(alert.AlertID),
			BasicRoomInfo:   alert.BasicDeviceInfo.BasicRoomInfo,
			RoomTags:        alert.RoomTags,
			AlertSeverities: []structs.AlertSeverity{alert.Severity},
			AlertTypes:      []structs.AlertType{alert.Type},
			AlertCategories: []structs.AlertCategory{alert.Category},
			SystemType:      alert.SystemType,
			Alerts: []structs.Alert{
				alert,
			},
		}

	} else {
		log.L.Errorf("Error: %v", err.Error())
		return
	}

	if issue.SystemType == "unknown" {
		issue, err = AddSystemTypeToIssue(issue)
		if err != nil {
			log.L.Errorf("Error getting system type for issue %v:%v", issue.RoomIssueID, err.Error())
		}
	}

	issue.CalculateAggregateInfo()
	// if roomAggregateChange {
	// 	calculateAggregateInfo(issue)
	// }

	err = alertcache.GetAlertCache("default").PutIssue(issue)
	if err != nil {
		log.L.Errorf("%v", "Couldn't save issue %v: %v", issue.RoomIssueID, err.Error())
		return
	}

	//run the alert change issue.
	a.runAlertActions(alert)

	if len(issue.RoomIssueResponses) == 0 {
		// auto-resolution rule
		if len(issue.ActiveAlertSeverities) == 0 {
			log.L.Debugf("Autoresolving Issue %v", issue.RoomIssueID)
			resInfo := structs.ResolutionInfo{
				Code:       "Auto Resolved",
				Notes:      "issue was auto resolved.",
				ResolvedAt: time.Now(),
			}

			err := a.resolveIssue(resInfo, issue.RoomIssueID, false, []string{})
			if err != nil {
				log.L.Errorf("Problem autoresolving issue %v: %v", issue.RoomIssueID, err.Error())
			}

			return
		} else if len(issue.ActiveAlertSeverities) < len(issue.AlertSeverities) {
			// its a partial resolution
			toClear := []structs.AlertSeverity{}
			for _, i := range issue.AlertSeverities {
				found := false
				for _, j := range issue.ActiveAlertSeverities {
					if i == j {
						found = true
						break
					}
				}
				if !found {
					toClear = append(toClear, i)
				}
			}

			log.L.Debugf("Room Issue %s, doing partial resolution for alert types %v", issue.RoomID, toClear)

			// for each alert severity to clear we're gonna do a partial resolution with those alerts
			resInfo := structs.ResolutionInfo{
				Code:       "Auto Resolved",
				Notes:      fmt.Sprintf("alerts for severity type(s) %v were auto resolved.", issue.AlertSeverities),
				ResolvedAt: time.Now(),
			}

			toResolve := []string{}
			for _, i := range issue.Alerts {
				for _, j := range toClear {
					if i.Severity == j {
						if i.Active {
							log.L.Errorf("Active alert found in partial clearing")
						} else {
							toResolve = append(toResolve, i.AlertID)
						}
					}
				}
			}

			//sumbit for partial resolution
			if len(toResolve) > 0 {
				log.L.Debugf("Room Issue %s, doing partial resolution for alerts %v", issue.RoomID, toResolve)

				err := a.resolveIssue(resInfo, issue.RoomIssueID, true, toResolve)
				if err != nil {
					log.L.Errorf("Problem doing a partial autoresolution issue %v: %v", issue.RoomIssueID, err.Error())
				}
			} else {
				log.L.Errorf("Room Issue %s, Alerts for serverity %v to be cleared, but no alerts found.  AlertSeverities: %v, ActiveAlertSeverities: %v",
					issue.RoomIssueID, toClear, issue.AlertSeverities, issue.ActiveAlertSeverities)
			}
		}
	}

	persist.GetElkAlertPersist().StoreIssue(issue, false, false)
	a.runIssueActions(issue)
	socket.GetManager().WriteToSockets(issue)
}

func (a *alertStore) runIssueActions(issue structs.RoomIssue) {
	if a.actionManager != nil {
		go func() {
			acts := actions.DefaultConfig().GetActionsByTrigger("issue-change")

			log.L.Debugf("Running %v issue change actions for issue %v", len(acts), issue.RoomIssueID)

			// a new context for the run of this action
			actx := actionctx.PutRoomIssue(context.Background(), issue)
			for i := range acts {
				a.actionManager.RunAction(actx, acts[i])
			}
		}()
	}
}

func (a *alertStore) runAlertActions(alert structs.Alert) {
	if a.actionManager != nil {
		go func() {
			acts := actions.DefaultConfig().GetActionsByTrigger("alert-change")

			log.L.Debugf("Running %v alert change actions for alert %v", len(acts), alert.AlertID)

			// a new context for the run of this action
			actx := actionctx.PutAlert(context.Background(), alert)
			for i := range acts {
				a.actionManager.RunAction(actx, acts[i])
			}
		}()
	}
}

func (a *alertStore) runAlertInitActions(alert structs.Alert) {
	if a.actionManager != nil {
		go func() {
			acts := actions.DefaultConfig().GetActionsByTrigger("alert-init")

			log.L.Debugf("Running %v alert init actions for alert %v", len(acts), alert.AlertID)

			// a new context for the run of this action
			actx := actionctx.PutAlert(context.Background(), alert)
			for i := range acts {
				a.actionManager.RunAction(actx, acts[i])
			}
		}()
	}
}

//NOT SAFE FOR CONCURRENT ACCESS. DO NOT USE OUTSIDE OF run()
func (a *alertStore) handleRequest(req issueRequest) {

	log.L.Infof("Handling request %+v", req)

	if req.All {
		toReturn, err := alertcache.GetAlertCache("default").GetAllIssues()

		for i := range toReturn {
			toReturn[i].Source = Cache
		}

		req.ResponseChan <- issueResponse{
			Error:      err,
			RoomIssues: toReturn,
		}

	} else {
		v, err := alertcache.GetAlertCache("default").GetIssue(req.RoomIssueID)
		v.Source = Cache

		req.ResponseChan <- issueResponse{
			Error:      err,
			RoomIssues: []structs.RoomIssue{v},
		}
	}

}

func combineIssues(n, o structs.RoomIssue) (structs.RoomIssue, bool) {
	changes := false

	if len(n.IssueTags) != len(o.IssueTags) || !structs.ContainsAllTags(n.IssueTags, o.IssueTags...) {
		o.IssueTags = n.IssueTags
		changes = true
	}

	log.L.Infof("Old Incidents: %v. New Incidents: %v", o.IncidentID, n.IncidentID)
	if len(n.IncidentID) > 0 && (len(n.IncidentID) != len(o.IncidentID) || !structs.ContainsAllTags(o.IncidentID, n.IncidentID...)) {
		o.IncidentID = n.IncidentID
		changes = true
	}

	if len(n.Notes) > 0 && n.Notes != o.Notes {
		o.NotesLog = append(o.NotesLog, n.Notes)
		o.Notes = n.Notes
		changes = true
	}

	if len(n.RoomIssueResponses) != len(o.RoomIssueResponses) {
		// TODO add logic for if one changed (i don't know if we actually need that)
		// so now changes is only true if the number of room issue responses changed, not if an invidual one changed
		// the old stuff doesn't work because it adds the lists together and then checks if they have different length
		o.RoomIssueResponses = n.RoomIssueResponses
		changes = true
	}

	for i := range o.RoomIssueResponses {
		// matching ones should be at the same index
		if o.RoomIssueResponses[i].HelpSentAt != n.RoomIssueResponses[i].HelpSentAt ||
			o.RoomIssueResponses[i].HelpArrivedAt != n.RoomIssueResponses[i].HelpArrivedAt ||
			!structs.HasAllPeople(o.RoomIssueResponses[i].Responders, n.RoomIssueResponses[i].Responders...) {
			o.RoomIssueResponses[i] = n.RoomIssueResponses[i]
			changes = true
		}
	}

	return o, changes
}
