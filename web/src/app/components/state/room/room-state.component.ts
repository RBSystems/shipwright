import { Component, OnInit, ViewChild, AfterViewInit } from "@angular/core";
import { DataService } from "src/app/services/data.service";
import { MatTableDataSource, MatPaginator, MatSort } from "@angular/material";
import { CombinedRoomState, StaticDevice } from "src/app/objects/static";
import { StringsService } from "src/app/services/strings.service";
import { PageEvent } from "@angular/material";

@Component({
  selector: "room-state",
  templateUrl: "./room-state.component.html",
  styleUrls: ["./room-state.component.scss"]
})
export class RoomStateComponent implements OnInit, AfterViewInit {
  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;
  dataSource = new MatTableDataSource<CombinedRoomState>([]);
  columns = ["type", "roomID", "alerts", "devices"];
  roomList: CombinedRoomState[] = [];
  filteredRoomList: CombinedRoomState[] = [];
  filterQueries: string[] = [];  
  length = 360;
  pageSize = 20;
  pageSizeOptions: number[] = [5, 10, 15, 20, 25, 30, 50, 100];

  // MatPaginator Output
  pageEvent: PageEvent;

  constructor(public data: DataService, public text: StringsService) {
    
  }

  ngOnInit() {
     
  }

  ngAfterViewInit() {
    if (this.data.finished) {
      console.log("Got the data");
      this.roomList = this.data.combinedRoomStateList;
      this.roomList.sort((a, b) => a.roomID.localeCompare(b.roomID))

      for (const room of this.roomList) {
        room.deviceStates.sort((a, b) => 
        {
          if (a.deviceType === b.deviceType) {
            return a.deviceName.localeCompare(b.deviceName);
          }

          return a.deviceType.localeCompare(b.deviceType);

        });
        console.log(room.deviceStates);
      }

      this.filteredRoomList = this.roomList;
      this.SetDataSource();
    } else {
      this.data.loaded.subscribe(() => {
        console.log("Subscribed to get the data");
        this.roomList = this.data.combinedRoomStateList;
        this.roomList.sort((a, b) => a.roomID.localeCompare(b.roomID))

        this.filteredRoomList = this.roomList;
        this.SetDataSource();
      });
    }
  }

  SetDataSource() {
    this.dataSource.data = this.filteredRoomList;
    setTimeout(() => {
      this.dataSource.paginator = this.paginator;
      this.dataSource.sort = this.sort;
      });
  }

  GetDeviceName(deviceID: string): string {
    const IDParts = deviceID.split("-");
    if (IDParts.length === 3) {
      return IDParts[2];
    }
    return deviceID;
  }

  GetDeviceTypes(rs: CombinedRoomState): string[] {
    let deviceTypes: string[] = [];
    if (rs != null && rs.deviceStates != null) {
      for (const ds of rs.deviceStates) {
        if (!deviceTypes.includes(ds.deviceType)) {
          deviceTypes.push(ds.deviceType);
        }
      }
      deviceTypes = deviceTypes.sort();
    }
    return deviceTypes;
  }

  GetStyle(ds: StaticDevice): string {
    let style = "default-chip";
    let deviceType = ds.deviceType;
    if (deviceType) deviceType = deviceType.toLowerCase();

    let power = ds.power;
    if (power) power = power.toLowerCase();

    if (deviceType === "display" || deviceType === "dmps") {
      if (power === "on") {
        style = "display-on";
        return style;
      }
      // style = "display-standby";
      return style;
    }
    if (deviceType === "via") {
      if (ds.currentUserCount > 0) {
        style = "display-on";
        return style;
      }
      // style = "display-standby";
      return style;
    }

    if (deviceType === "microphone") {
      if (ds.batteryChargeMinutes >= 300) {
        style = "mic-good";
        return style;
      }
      if (300 > ds.batteryChargeMinutes && ds.batteryChargeMinutes >= 90) {
        style = "mic-okay";
        return style;
      }
      if (90 > ds.batteryChargeMinutes && ds.batteryChargeMinutes >= 0) {
        style = "mic-bad";
        return style;
      }
    }

    return style;
  }

  HasBatteryCharge(ds: StaticDevice): boolean {
    return ds.batteryChargeMinutes != null;
  }

  checkRoom(room, query) {
    for (const item of query) {
      if (room.roomID.toLowerCase().includes(item.toLowerCase())) {
        continue;
      }
      if (room.staticRoom.systemType.length > 0) {
        if (
          room.staticRoom.systemType[0]
            .toLowerCase()
            .includes(item.toLowerCase())
        ) {
          continue;
        }
      }
      if (room.deviceStates != null) {
        let result = false;
        for (const device of room.deviceStates) {
          if (
            device.deviceID.toLowerCase().includes(item.toLowerCase()) &&
            !this.filteredRoomList.includes(room)
          ) {
            result = true;
            break;
          }
        }
        if (result) {
          continue;
        }
      }
      if (item === "Has Alerts") {
        console.log("has alerts is suposed to run");
        if (this.activeAlerts(room)) {
          console.log("returned true");
          continue;
        }
      }
      if (item === "Low Mic") {
        if (this.lowMic(room)) {
          continue;
        }
      }
      if (item === "Warn Mic") {
        if (this.warnMic(room)) {
          continue;
        }
      }
      if (item === "System In Use") {
        if (this.inUse(room)) {
          continue;
        }
      }
      return false;
    }
    return true;
  }

  FilterRooms() {

    this.filteredRoomList = [];
    if (this.filterQueries.length === 0) {
      this.filteredRoomList = this.roomList;
      this.SetDataSource();
      return;
    }

    for (const room of this.roomList) {
      if (this.checkRoom(room, this.filterQueries)) {
        this.filteredRoomList.push(room);
      }
    }

    this.SetDataSource();
  }

  activeAlerts(room?: CombinedRoomState) {
    if (!this.filterQueries.includes("Has Alerts")) {
      this.filterQueries.push("Has Alerts");
      this.FilterRooms();
    }
    if (room) {
      if (room.activeAlertCount > 0) {
        return true;
      }
    return false;
    }

  }

  lowMic(room?: CombinedRoomState) {
    if (!this.filterQueries.includes("Low Mic")) {
      this.filterQueries.push("Low Mic");
      this.FilterRooms();
    }
    if (room) {
      if (room.deviceStates != null) {
        for (const device of room.deviceStates) {
          if (
            device.deviceType === "microphone" &&
            90 > device.batteryChargeMinutes &&
            device.batteryChargeMinutes >= 0 &&
            !this.filteredRoomList.includes(room)
          ) {
            return true;
          }
        }
      }
      return false;
    }
  }

  warnMic(room?: CombinedRoomState) {
    if (!this.filterQueries.includes("Warn Mic")) {
      this.filterQueries.push("Warn Mic");
      this.FilterRooms();
    }
    if (room) {
      if (room.deviceStates != null) {
        for (const device of room.deviceStates) {
          if (
            device.deviceType === "microphone" &&
            300 > device.batteryChargeMinutes &&
            device.batteryChargeMinutes >= 90 &&
            !this.filteredRoomList.includes(room)
          ) {
            return true;
          }
        }
      }
      return false;
    }
  }

  inUse(room?: CombinedRoomState) {
    if (!this.filterQueries.includes("System In Use")) {
      this.filterQueries.push("System In Use");
      this.FilterRooms();
    }
    if (room) {
      if (room.deviceStates != null) {
        for (const device of room.deviceStates) {
          let deviceType = device.deviceType;
          if (deviceType) {deviceType = deviceType.toLowerCase(); }

          let power = device.power;
          if (power) {power = power.toLowerCase(); }

          if (deviceType === "display" || deviceType === "dmps") {
            if (
              power === "on" &&
              !this.filteredRoomList.includes(room)
            ) {
              return true;
            }
          }
        }
      }
      return false;
    }
  }
}
