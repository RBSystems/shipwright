import {
  JsonObject,
  JsonProperty,
  JsonConverter,
  JsonCustomConvert,
  Any,
} from "json2typescript";
import { Person } from "src/app/objects/database";

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
  resolvedAt: Date = new Date("1970-01-01T00:00:00.000Z");
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

  @JsonProperty("system-type", String, true)
  systemType: string = undefined;

  @JsonProperty("start-time", DateConverter, true)
  startTime: Date = undefined;

  @JsonProperty("end-time", DateConverter, true)
  endTime: Date = undefined;

  @JsonProperty("update-time", DateConverter, true)
  updateTime: Date = undefined;

  @JsonProperty("active", Boolean, true)
  active: boolean = undefined;

  @JsonProperty("alert-tags", [String], true)
  alertTags: string[] = Array<string>();

  @JsonProperty("device-tags", [String], true)
  deviceTags: string[] = Array<string>();

  @JsonProperty("requester", String, true)
  requester: string = undefined;

  @JsonProperty("manual-resolve", Boolean, true)
  manualResolve: boolean = undefined;
}

@JsonObject("RoomIssueResponse")
export class RoomIssueResponse {
  @JsonProperty("responders", [String], true)
  responders: Person[] = Array<Person>();

  @JsonProperty("help-sent-at", DateConverter, true)
  helpSentAt: Date = undefined;

  @JsonProperty("help-arrived-at", DateConverter, true)
  helpArrivedAt: Date = undefined;

  SentIsZero(): boolean {
    if (this.helpSentAt === undefined) {
      return true;
    }

    const zero = "0001-01-01T00:00:00.000Z";

    return this.helpSentAt.toISOString() === zero;
  }

  ArrivedIsZero(): boolean {
    if (this.helpArrivedAt === undefined) {
      return true;
    }
    const zero = "0001-01-01T00:00:00.000Z";

    return this.helpArrivedAt.toISOString() === zero;
  }
}

@JsonObject("RoomIssue")
export class RoomIssue {
  @JsonProperty("id", String, true)
  issueID: string = undefined;

  @JsonProperty("buildingID", String, true)
  buildingID: string = undefined;

  @JsonProperty("roomID", String, true)
  roomID: string = undefined;

  @JsonProperty("alert-severities", [String], true)
  alertSeverities: string[] = undefined;

  @JsonProperty("room-tags", [String], true)
  roomTags: string[] = Array<string>();

  @JsonProperty("issue-tags", [String], true)
  issueTags: string[] = Array<string>();

  @JsonProperty("alert-types", [String], true)
  alertTypes: string[] = Array<string>();

  @JsonProperty("alert-devices", [String], true)
  alertDevices: string[] = Array<string>();

  @JsonProperty("alert-categories", [String], true)
  alertCategories: string[] = Array<string>();

  @JsonProperty("alert-count", Number, true)
  alertCount: number = undefined;

  @JsonProperty("active-alert-types", [String], true)
  activeAlertTypes: string[] = Array<string>();

  @JsonProperty("active-alert-categories", [String], true)
  activeAlertCategories: string[] = Array<string>();

  @JsonProperty("active-alert-devices", [String], true)
  activeAlertDevices: string[] = Array<string>();

  @JsonProperty("active-alert-severities", [String], true)
  activeAlertSeverities: string[] = Array<string>();

  @JsonProperty("active-alert-count", Number, true)
  activeAlertCount: number = undefined;

  @JsonProperty("incident-id", [String], true)
  incidentID: string[] = Array<string>();

  @JsonProperty("system-type", String)
  systemType: string = undefined;

  @JsonProperty("notes", String, true)
  notes: string = undefined;

  @JsonProperty("notes-log", [String], true)
  notesLog: string[] = Array<string>();

  @JsonProperty("responses", [RoomIssueResponse], true)
  roomIssueResponses: RoomIssueResponse[] = Array<RoomIssueResponse>();

  @JsonProperty("resolved", Boolean, true)
  resolved: boolean = undefined;

  @JsonProperty("resolution-info", ResolutionInfo, true)
  resolutionInfo: ResolutionInfo = new ResolutionInfo();

  @JsonProperty("alerts", [Alert], true)
  alerts: Alert[] = Array<Alert>();

  ResolvedAtIsZero(): boolean {
    if (this.resolutionInfo.resolvedAt === undefined) {
      return true;
    }
    const zero = "0001-01-01T00:00:00.000Z";

    return this.resolutionInfo.resolvedAt.toISOString() === zero;
  }
}

@JsonObject("ClassHalfHourBlock")
export class ClassHalfHourBlock {
  @JsonProperty("block-start", String, true)
  blockStart: string = undefined;

  @JsonProperty("class-name", String, true)
  className: string = undefined;

  @JsonProperty("class-time", String, true)
  classTime: string = undefined;

  @JsonProperty("teacher", Person, true)
  teacher: Person = undefined;

  @JsonProperty("days", String, true)
  days: string = undefined;

  @JsonProperty("class-start-time", DateConverter, true)
  classStartTime: Date = undefined;

  @JsonProperty("class-end-time", DateConverter, true)
  classEndTime: Date = undefined;
}
