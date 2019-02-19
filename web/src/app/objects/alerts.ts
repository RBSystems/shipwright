import { JsonObject, JsonProperty, JsonConverter, JsonCustomConvert, Any } from "json2typescript";

@JsonConverter
class DateConverter implements JsonCustomConvert<Date> {
    serialize(date: Date): any {
        function pad(n) {
        return n < 10 ? "0" + n : n;
        }
    
        return (
        date.getUTCFullYear() +
        "-" +
        pad(date.getUTCMonth() + 1) +
        "-" +
        pad(date.getUTCDate()) +
        "T" +
        pad(date.getUTCHours()) +
        ":" +
        pad(date.getUTCMinutes()) +
        ":" +
        pad(date.getUTCSeconds()) +
        "Z"
        );
    }
    
    deserialize(date: any): Date {
        if (date == null) {
        return undefined;
        }
    
        return new Date(date);
    }
}

@JsonObject("ResolutionInfo")
export class ResolutionInfo {
    @JsonProperty("resolution-code", String, true)
    code: string = undefined;
    
    @JsonProperty("notes", String, true)
    notes: string = undefined;
    
    @JsonProperty("resolved-at", DateConverter, true)
    resolvedAt: Date = new Date(0);
}

@JsonObject("Alert")
export class Alert {
    @JsonProperty("id", String, true)
    alertID: string = undefined;
    
    @JsonProperty("buildingID", String, true)
    buildingID: string = undefined;
    
    @JsonProperty("roomID", String, true)
    roomID: string = undefined;
    
    @JsonProperty("deviceID", String, true)
    deviceID: string = undefined;
    
    @JsonProperty("type", String, true)
    type: string = undefined;
    
    @JsonProperty("category", String, true)
    category: string = undefined;

    @JsonProperty("severity", String, true)
    severity: string = undefined;
    
    @JsonProperty("message", String, true)
    message: string = undefined;
    
    @JsonProperty("message-log", [String], true)
    messageLog: string[] = Array<string>();
    
    @JsonProperty("data", Any, true)
    data: any = undefined;

    @JsonProperty("system-type", String)
    systemType: string = undefined;
    
    @JsonProperty("start-time", DateConverter, true)
    startTime: Date = new Date(0);
    
    @JsonProperty("end-time", DateConverter, true)
    endTime: Date = new Date(0);
    
    @JsonProperty("update-time", DateConverter, true)
    updateTime: Date = new Date(0);
    
    @JsonProperty("active", Boolean, true)
    active: boolean = undefined;
    
    @JsonProperty("alert-tags", [String], true)
    alertTags: string[] = Array<string>();
    
    @JsonProperty("device-tags", [String], true)
    deviceTags: string[] = Array<string>();
}

@JsonObject("RoomIssue")
export class RoomIssue {
    @JsonProperty("id", String, true)
    issueID: string = undefined;
    
    @JsonProperty("buildingID", String, true)
    buildingID: string = undefined;
    
    @JsonProperty("roomID", String, true)
    roomID: string = undefined;

    @JsonProperty("severity", String, true)
    severity: string = undefined;

    @JsonProperty("room-tags", [String], true)
    roomTags: string[] = Array<string>();

    @JsonProperty("issue-tags", [String], true)
    issueTags: string[] = Array<string>();

    @JsonProperty("alert-types", [String], true)
    alertTypes: string[] = Array<string>();

    @JsonProperty("alert-categories", [String], true)
    alertCategories: string[] = Array<string>();

    @JsonProperty("incident-id", String, true)
    incidentID: string = undefined;

    @JsonProperty("system-type", String)
    systemType: string = undefined;

    @JsonProperty("notes", String, true)
    message: string = undefined;
    
    @JsonProperty("notes-log", [String], true)
    messageLog: string[] = Array<string>();

    @JsonProperty("responders", [String], true)
    responders: string[] = Array<string>();
    
    @JsonProperty("help-sent-at", DateConverter, true)
    helpSentAt: Date = new Date(0);
    
    @JsonProperty("help-arrived-at", DateConverter, true)
    helpArrivedAt: Date = new Date(0);

    @JsonProperty("resolved", Boolean, true)
    resolved: boolean = undefined;
    
    @JsonProperty("resolution-info", ResolutionInfo, true)
    resolutionInfo: ResolutionInfo = undefined;

    @JsonProperty("alerts", [Alert], true)
    alerts: Alert[] = Array<Alert>();

    SentIsZero(): boolean {
        let zero = "0001-01-01T00:00:00.000Z";

        return this.helpSentAt.toISOString() === zero;
    }

    ArrivedIsZero(): boolean {
        let zero = "0001-01-01T00:00:00.000Z";

        return this.helpArrivedAt.toISOString() === zero;
    }

    GetAlerts(severity?: string): Alert[] {
        let toReturn: Alert[] = [];

        this.alerts.forEach((alert, id) => {
            if(severity != null) {
                if(alert.severity == severity) {
                    toReturn.push(alert);
                }
            } else {
                toReturn.push(alert);
            }
        })
        
        return toReturn;
    }

    ActiveAlertCount(): number {
        let count = 0;

        for(let alert of this.GetAlerts()) {
            if(alert.active) {
                count++
            }
        }

        return count
    }

    TotalAlertCount(): number {
        let count = this.GetAlerts().length
        return count
    }
}

export const AllAlerts = "all-alerts";
export const CritAlerts = "critical"
export const WarnAlerts = "warning"
export const Battery = "battery"

