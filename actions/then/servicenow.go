package then

import (
	"context"
	"sync"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/servicenow"
)

//this may or may not be right...
type Func func(ctx context.Context, with []byte) *nerr.E

var (
	thens struct {
		sync.RWMutex
		m map[string]Func
	}
)

//would this need to return the incident? or what are we doing with it/where are we
//storing incident info?

//TODO: create only one then statement (if there is no incident number Then create incident)
//If there is an incident number and it is not resolved yet, then modify incident
//if there is an incident number AND it is resolved, then close the incident.
func CreateIncident(ctx context.Context, with []byte) *nerr.E {
	alert, ok := alertctx.GetAlert(ctx)
	if !ok {
		log.L.Errorf("Failed to get Alert")
		return "", nerr.Create("Must have alert to create incident", "")
	}

	//pass alert and call the create incident function is ServiceNow
	if alert.IncidentID == nil {
		incident, err := servicenow.CreateIncident(alert)
		if err != nil {
			log.L.Errorf("Failed to create incident")
			return "", nerr.Translate(err).Add("Incident was not created in servicenow")
		}
		alert.IncidentID = incident.Number
		AddAlert(alert)
	}
	if alert.IncidentID != nil && alert.Resolved == false {
		incident, err := servicenow.ModifyIncident(alert)
		if err != nil {
			log.L.Errorf("Failed to Modify incident")
			return "", nerr.Translate(err).Add("Incident was not modified in servicenow")
		}
	}

	if alert.IncidentID != nil && alert.Resolved == true {
		incident, err := servicenow.CloseIncident(alert)
		if err != nil {
			log.L.Errorf("Failed to close incident")
			return "", nerr.Translate(err).Add("Incident was not closed in servicenow")
		}
	}
	return nil
}
