import { Component, Injectable } from '@angular/core';
import { MatDialog } from '@angular/material';
import { SettingsModalComponent } from '../modals/settings/settings.component';
import { RespondModalComponent } from '../modals/respond/respond.component';
import { BuildingModalComponent } from '../modals/buildingmodal/buildingmodal.component';
import { RoomModalComponent } from '../modals/roommodal/roommodal.component';
import { DeviceModalComponent } from '../modals/devicemodal/devicemodal.component';
import { NotifyModalComponent } from '../modals/notify/notify.component';
import { PresetModalComponent } from '../modals/presetmodal/presetmodal.component';
import { IconModalComponent } from '../modals/iconmodal/iconmodal.component';
import { RoomIssue, Alert } from '../objects/alerts';
import { Building, Room, Device, DBResponse, Preset, Panel, UIConfig } from '../objects/database';
import { APIService } from './api.service';

@Injectable({
  providedIn: 'root'
})
export class ModalService {

  constructor(private dialog: MatDialog, private api: APIService) { }

  OpenSettingsModal() {
    this.dialog.open(SettingsModalComponent);
  }

  OpenRespondModal(issue: RoomIssue, selected: Alert[]) {
    this.dialog.open(RespondModalComponent, {data: {issue: issue, selected: selected}}).afterClosed().subscribe((issue) => {
      if(issue != null) {
        this.api.ResolveIssue(issue)
      }
    })
  }

  OpenBuildingModal(building: Building) {
    this.dialog.open(BuildingModalComponent, {data: building}).afterClosed().subscribe((resp) => {
      if (resp != null){
        this.OpenNotifyModal(resp)
      }
    })
  }

  OpenRoomModal(room: Room) {
    this.dialog.open(RoomModalComponent, {data: room}).afterClosed().subscribe((resp) => {
      if (resp != null){
        this.OpenNotifyModal(resp)
      }
    })
  }

  OpenDeviceModal(device: Device) {
    this.dialog.open(DeviceModalComponent, {data: device});
  }

  OpenNotifyModal(resp: DBResponse) {
    this.dialog.open(NotifyModalComponent, {data: resp});
  }

  OpenPresetModal(preset: Preset, cps: Panel[], config: UIConfig) {
    this.dialog.open(PresetModalComponent, {data: {preset: preset, currentPanels: cps, config: config}})
  }

  OpenIconModal(caller: any) {
    this.dialog.open(IconModalComponent).afterClosed().subscribe(result => {
      if(result != null) {
        caller.icon = result;
      }
    });
  }
}
