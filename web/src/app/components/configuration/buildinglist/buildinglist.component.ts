import { Component, OnInit } from "@angular/core";
import { StringsService } from "src/app/services/strings.service";
import { DataService } from "src/app/services/data.service";
import { Building } from "src/app/objects/database";
import { ModalService } from "src/app/services/modal.service";

@Component({
  selector: "building-list",
  templateUrl: "./buildinglist.component.html",
  styleUrls: ["./buildinglist.component.scss"]
})
export class BuildingListComponent implements OnInit {
  buildingList: Building[] = [];
  filteredBuildings: Building[] = [];
  filterQueries: string[] = [];

  constructor(public text: StringsService, public data: DataService, public modal: ModalService) {
    if (this.data.finished) {
      this.buildingList = this.data.allBuildings;
      this.filteredBuildings = this.buildingList;
    } else {
      this.data.loaded.subscribe(() => {
        this.buildingList = this.data.allBuildings;
        this.filteredBuildings = this.buildingList;
      });
    }
  }

  ngOnInit() {
  }

  GoBack() {
    window.history.back();
  }

  FilterBuildings() {
    this.filteredBuildings = [];

    if (this.filterQueries.length === 0) {
      this.filteredBuildings = this.buildingList;
      return;
    }

    for (const building of this.buildingList) {
        for (const query of this.filterQueries) {
          if (building.id.toLowerCase().includes(query.toLowerCase()) && !this.filteredBuildings.includes(building)) {
            this.filteredBuildings.push(building);
          }
          if (building.description.toLowerCase().includes(query.toLowerCase()) && !this.filteredBuildings.includes(building)) {
            this.filteredBuildings.push(building);
          }
          for (const tag of building.tags) {
            if (tag.toLowerCase().includes(query.toLowerCase()) && !this.filteredBuildings.includes(building)) {
              this.filteredBuildings.push(building);
            }
          }
      }
    }
  }

  CreateNewBuilding() {
    const newBuilding = new Building();
    newBuilding.isNew = true;
    this.modal.OpenBuildingModal(newBuilding);
    this.modal.buildingDone.subscribe((b) => {
      this.filteredBuildings.push(b);
      this.filteredBuildings.sort(this.text.SortAlphaNumByID);
    });
  }
}
