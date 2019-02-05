import { Injectable, EventEmitter } from '@angular/core';
import { JsonConvert, Any } from 'json2typescript';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Building, DBResponse, Room, RoomConfiguration, Device, DeviceType, Role, UIConfig, BuildingStatus, RoomStatus, Template, MetricsResponse } from '../objects';

@Injectable({
  providedIn: 'root'
})
export class APIService {
  public theme = "default";
  public themeSwitched: EventEmitter<string[]>;

  private converter: JsonConvert;
  private urlParams: URLSearchParams;

  private headers: HttpHeaders;

  constructor(private http: HttpClient) {
    this.themeSwitched = new EventEmitter<string[]>();
    this.converter = new JsonConvert();
    this.converter.ignorePrimitiveChecks = false;

    this.urlParams = new URLSearchParams(window.location.search);
    if (this.urlParams.has("theme")) {
      this.theme = this.urlParams.get("theme");
    }

    this.headers = new HttpHeaders(
      {'Content-Type': 'application/json',
      'Access-Control-Allow-Origin': '*'}
    );
  }

  public refresh() {
    window.location.reload(true);
  }

  public switchTheme(name: string) {
    let oldTheme = this.theme+"-theme"
    let newTheme = name+"-theme"

    console.log("switching theme to ", name);

    this.theme = name;
    this.urlParams.set("theme", name);
    
    this.themeSwitched.emit([oldTheme, newTheme]);

    window.history.replaceState(
      null,
      "Shipwright",
      window.location.pathname + "?" + this.urlParams.toString()
    );    
  }

  // Building Functions
  public async AddBuilding(toAdd: Building) {
    try {
      const data = await this.http.post("buildings/"+toAdd.id, toAdd, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error adding the building " + toAdd.id + ": " + e);
    }
  }

  public async AddBuildings(toAdd: Building[]) {
    try {
      const data = await this.http.post("buildings", toAdd, {headers: this.headers})
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error adding multiple buildings: " + e);
    }
  }

  public async GetBuilding(buildingID: string) {
    try {
      const data = await this.http.get("buildings/"+buildingID, {headers: this.headers}).toPromise();
      const building = this.converter.deserialize(data, Building);

      return building;
    } catch (e) {
      throw new Error("error getting the building " + buildingID + ": " + e);
    }
  }

  public async GetAllBuildings() {
    try {
      const data = await this.http.get("buildings", {headers: this.headers}).toPromise();
      const buildings = this.converter.deserialize(data, Building);

      return buildings;
    } catch (e) {
      throw new Error("error getting all buildings: " + e);
    }
  }

  public async UpdateBuilding(idToUpdate: string, toUpdate: Building) {
    try {
      console.log(idToUpdate);
      console.log(toUpdate);
      const data = await this.http.put("buildings/"+idToUpdate+"/update", toUpdate, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error updating the building " + idToUpdate + ": " + e);
    }
  }

  public async UpdateBuildings(toUpdate: Building[]) {
    try {
      const data = await this.http.put("buildings/update", toUpdate, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error updating multiple buildings: " + e);
    }
  }

  public async DeleteBuilding(buildingID: string) {
    try {
      const data = await this.http.get("buildings/"+buildingID+"/delete", {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error deleting the building " + buildingID + ": " + e);
    }
  }

  // Room Functions
  public async AddRoom(toAdd: Room) {
    try {
      const data = await this.http.post("rooms/"+toAdd.id, toAdd, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error adding the room " + toAdd.id + ": " + e);
    }
  }

  public async AddRooms(toAdd: Room[]) {
    try {
      const data = await this.http.post("rooms", toAdd, {headers: this.headers})
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error adding multiple rooms: " + e);
    }
  }

  public async GetRoom(roomID: string) {
    try {
      const data = await this.http.get("rooms/"+roomID, {headers: this.headers}).toPromise();
      const room = this.converter.deserialize(data, Room);

      return room;
    } catch (e) {
      throw new Error("error getting the room " + roomID + ": " + e);
    }
  }

  public async GetAllRooms() {
    try {
      const data = await this.http.get("rooms", {headers: this.headers}).toPromise();
      const rooms = this.converter.deserialize(data, Room);

      return rooms;
    } catch (e) {
      throw new Error("error getting all rooms: " + e);
    }
  }

  public async GetRoomsByBuilding(buildingID: string) {
    try {
      const data = await this.http.get("buildings/"+buildingID+"/rooms", {headers: this.headers}).toPromise();
      const rooms = this.converter.deserialize(data, Room);

      return rooms;
    } catch (e) {
      throw new Error("error getting all rooms in the building " + buildingID + ": " + e);
    }
  }

  public async UpdateRoom(idToUpdate: string, toUpdate: Room) {
    try {
      const data = await this.http.put("rooms/"+idToUpdate+"/update", toUpdate, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error updating the room " + idToUpdate + ": " + e);
    }
  }

  public async UpdateRooms(toUpdate: Room[]) {
    try {
      const data = await this.http.put("rooms/update", toUpdate, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error updating multiple rooms: " + e);
    }
  }

  public async DeleteRoom(roomID: string) {
    try {
      const data = await this.http.get("rooms/"+roomID+"/delete", {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error deleting the room " + roomID + ": " + e);
    }
  }

  public async GetRoomConfigurations() {
    try {
      const data = await this.http.get("rooms/configurations", {headers: this.headers}).toPromise();
      const roomConfigs = this.converter.deserialize(data, RoomConfiguration);

      return roomConfigs;
    } catch (e) {
      throw new Error("error getting all room configurations: " + e);
    }
  }

  public async GetRoomDesignations() {
    try {
      const data = await this.http.get("rooms/designations", {headers: this.headers}).toPromise();

      return data;
    } catch (e) {
      throw new Error("error getting all room designations: " + e);
    }
  }

  // Device Functions
  public async AddDevice(toAdd: Device) {
    try {
      const data = await this.http.post("devices/"+toAdd.id, toAdd, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error adding the device " + toAdd.id + ": " + e);
    }
  }

  public async AddDevices(toAdd: Device[]) {
    try {
      const data = await this.http.post("devices", toAdd, {headers: this.headers})
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error adding multiple devices: " + e);
    }
  }

  public async GetDevice(deviceID: string) {
    try {
      const data = await this.http.get("devices/"+deviceID, {headers: this.headers}).toPromise();
      const device = this.converter.deserialize(data, Device);

      return device;
    } catch (e) {
      throw new Error("error getting the device " + deviceID + ": " + e);
    }
  }

  public async GetAllDevices() {
    try {
      const data = await this.http.get("devices", {headers: this.headers}).toPromise();
      const devices = this.converter.deserialize(data, Device);

      return devices;
    } catch (e) {
      throw new Error("error getting all devices: " + e);
    }
  }

  public async GetDevicesByRoom(roomID: string) {
    try {
      const data = await this.http.get("rooms/"+roomID+"/devices", {headers: this.headers}).toPromise();
      const devices = this.converter.deserialize(data, Device);

      return devices;
    } catch (e) {
      throw new Error("error getting all devices in the room " + roomID + ": " + e);
    }
  }

  public async GetDevicesByRoomAndRole(roomID: string, roleID: string) {
    try {
      const data = await this.http.get("rooms/"+roomID+"/devices/roles/"+roleID, {headers: this.headers}).toPromise();
      const devices = this.converter.deserialize(data, Device);

      return devices;
    } catch (e) {
      throw new Error("error getting all devices in the room " + roomID + " that have the role " + roleID + ": " + e);
    }
  }

  public async GetDevicesByTypeAndRole(typeID: string, roleID: string) {
    try {
      const data = await this.http.get("devices/types/"+typeID+"/roles/"+roleID, {headers: this.headers}).toPromise();
      const devices = this.converter.deserialize(data, Device);

      return devices;
    } catch (e) {
      throw new Error("error getting all devices of the type " + typeID + " that have the role " + roleID + ": " + e);
    }
  }

  public async UpdateDevice(idToUpdate: string, toUpdate: Device) {
    try {
      const data = await this.http.put("devices/"+idToUpdate+"/update", toUpdate, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error updating the device " + idToUpdate + ": " + e);
    }
  }

  public async UpdateDevices(toUpdate: Device[]) {
    try {
      const data = await this.http.put("devices/update", toUpdate, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error updating multiple devices: " + e);
    }
  }

  public async DeleteDevice(deviceID: string) {
    try {
      const data = await this.http.get("devices/"+deviceID+"/delete", {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error deleting the device " + deviceID + ": " + e);
    }
  }

  public async GetDeviceTypes() {
    try {
      const data = await this.http.get("devices/types", {headers: this.headers}).toPromise();
      const deviceTypes = this.converter.deserialize(data, DeviceType);

      return deviceTypes;
    } catch (e) {
      throw new Error("error getting all device types: " + e);
    }
  }

  public async GetDeviceRoles() {
    try {
      const data = await this.http.get("devices/roles", {headers: this.headers}).toPromise();
      const deviceRoles = this.converter.deserialize(data, Role);

      return deviceRoles;
    } catch (e) {
      throw new Error("error getting all device roles: " + e);
    }
  }

  // UIConfig Functions
  public async AddUIConfig(toAdd: UIConfig) {
    try {
      const data = await this.http.post("uiconfigs/"+toAdd.id, toAdd, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error adding the uiconfig " + toAdd.id + ": " + e);
    }
  }

  public async AddUIConfigs(toAdd: UIConfig[]) {
    try {
      const data = await this.http.post("uiconfigs", toAdd, {headers: this.headers})
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error adding multiple uiconfigs: " + e);
    }
  }

  public async GetUIConfig(configID: string) {
    try {
      const data = await this.http.get("uiconfigs/"+configID, {headers: this.headers}).toPromise();
      const uiconfig = this.converter.deserialize(data, UIConfig);

      return uiconfig;
    } catch (e) {
      throw new Error("error getting the uiconfig " + configID + ": " + e);
    }
  }

  public async GetAllUIConfigs() {
    try {
      const data = await this.http.get("uiconfigs", {headers: this.headers}).toPromise();
      const uiconfigs = this.converter.deserialize(data, UIConfig);

      return uiconfigs;
    } catch (e) {
      throw new Error("error getting all uiconfigs: " + e);
    }
  }

  public async UpdateUIConfig(idToUpdate: string, toUpdate: UIConfig) {
    try {
      const data = await this.http.put("uiconfigs/"+idToUpdate+"/update", toUpdate, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error updating the uiconfig " + idToUpdate + ": " + e);
    }
  }

  public async UpdateUIConfigs(toUpdate: UIConfig[]) {
    try {
      const data = await this.http.put("uiconfigs/update", toUpdate, {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error updating multiple uiconfigs: " + e);
    }
  }

  public async DeleteUIConfig(configID: string) {
    try {
      const data = await this.http.get("uiconfigs/"+configID+"/delete", {headers: this.headers}).toPromise();
      const response = this.converter.deserialize(data, DBResponse);

      return response;
    } catch (e) {
      throw new Error("error deleting the uiconfig " + configID + ": " + e);
    }
  }

  // Options Functions
  public async GetIcons() {
    try {
      const data = await this.http.get("options/icons", {headers: this.headers}).toPromise();

      return data;
    } catch (e) {
      throw new Error("error getting the list of icons: " + e);
    }
  }

  public async GetTemplates() {
    try {
      const data = await this.http.get("options/templates", {headers: this.headers}).toPromise();
      const templates = this.converter.deserialize(data, Template);

      return templates;
    } catch (e) {
      throw new Error("error getting the list of templates: " + e);
    }
  }

  // Authentication Functions
  public async GetCurrentUsername() {
    try {
      const data = await this.http.get("users/current/username", {headers: this.headers}).toPromise();

      return data;
    } catch (e) {
      throw new Error("error getting the current user's username: " + e);
    }
  }

  public async GetUserPermissions() {
    try {
      const data = await this.http.get("users/current/permissions", {headers: this.headers}).toPromise();
      const permissions = this.converter.deserialize(data, Any);

      return permissions;
    } catch (e) {
      throw new Error("error getting the current user's permissions: " + e);
    }
  }
}