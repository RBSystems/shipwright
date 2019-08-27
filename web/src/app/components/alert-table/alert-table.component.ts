import { Component, OnInit, Input, ViewChild } from "@angular/core";
import { APIService } from "src/app/services/api.service";
import { COMMA, ENTER } from "@angular/cdk/keycodes";
import { RoomIssue } from "src/app/objects/alerts";
import {
  animate,
  state,
  style,
  transition,
  trigger
} from "@angular/animations";
import { Router } from "@angular/router";
import {
  MatPaginator,
  MatSort,
  MatTableDataSource,
  MatChipInputEvent
} from "@angular/material";

import { Filter, FilterType } from "../state/device/device-state.component";
import { DashpanelTypes } from "../dashpanel/idashpanel";

@Component({
  selector: "alert-table",
  templateUrl: "./alert-table.component.html",
  styleUrls: ["./alert-table.component.scss"],
  animations: [
    trigger("detailExpand", [
      state("collapsed", style({ height: "0px", minHeight: "0" })),
      state("expanded", style({ height: "*" })),
      transition(
        "expanded <=> collapsed",
        animate("225ms cubic-bezier(0.4, 0.0, 0.2, 1)")
      )
    ])
  ]
})
export class AlertTableComponent implements OnInit {
  readonly separatorKeyCodes: number[] = [ENTER, COMMA]; // delimate filters with these keys
  readonly filterType: typeof FilterType = FilterType; // so the component can use them

  @Input() singleRoom = false;
  @Input() roomID: string;
  @Input() chosenType: DashpanelTypes;

  issues: RoomIssue[];
  expandedIssue: RoomIssue | null;

  dataSource: MatTableDataSource<RoomIssue>;
  filters: Filter[] = [];

  issueCols = ["systemType", "roomID", "count", "age", "lastNote"];

  alertCols = ["deviceID", "type", "message", "startTime"];

  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;

  constructor(public api: APIService, public router: Router) {}

  ngOnInit() {
    this.api.GetAllIssues().then(answer => {
      this.issues = answer as RoomIssue[];
      console.log("issues", this.issues);
      this.dataSource = new MatTableDataSource(this.issues);

      // because of the ngIf on the parent container
      setTimeout(() => {
        this.dataSource.paginator = this.paginator;
        this.dataSource.sort = this.sort;
        this.dataSource.sortingDataAccessor = (
          data: RoomIssue,
          headerID: string
        ) => {
          switch (headerID) {
            case "age":
              return data.oldestActiveAlert.startTime;
            case "count":
              return data.activeAlertCount;
            default:
              return data[headerID];
          }
        };
      });

      this.dataSource.filterPredicate = (
        data: RoomIssue,
        filter: string
      ): boolean => {
        for (const f of this.filters) {
          if (!f.filter(data)) {
            return false;
          }
        }

        return true;
      };
    });
  }

  addFilter(ftype: FilterType, key: string, val: string) {
    const f = new Filter(ftype, key, val);
    this.filters.push(f);

    this.forceFilter();
  }

  removeFilter(filter: Filter) {
    const index = this.filters.indexOf(filter);

    if (index >= 0) {
      this.filters.splice(index, 1);
    }

    this.forceFilter();
  }

  addChip(event: MatChipInputEvent): void {
    const value = event.value.trim();

    let split = value.split(/:(.*)/);
    split = split.filter(s => s); // filter out blank ones
    split = split.map(s => s.trim()); // trim each string

    if (split.length === 2 && this.issueCols.includes(split[0])) {
      this.addFilter(FilterType.For, split[0], split[1]);
    } else {
      this.addFilter(FilterType.General, undefined, value);
    }

    if (event.input) {
      event.input.value = ""; // reset the input
    }
  }

  public forceFilter() {
    if (!this.dataSource.filter && this.filters.length > 0) {
      this.dataSource.filter = "◬";
    }

    this.dataSource._filterData(this.dataSource.data);
    this.dataSource._updateChangeSubscription();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  goToAlerts(roomID: string) {
    this.router.navigate(["/campus/" + roomID + "/tab/2"]);
  }

  getTotalAlertCount() {
    const count = 0;

    for (const issue of this.issues) {
      if (issue.alertCount !== undefined) {
        count += issue.alertCount;
      }
    }
    return count;
  }

  getTotalActiveAlertCount() {
    let count = 0;

    for (const issue of this.issues) {
      if (issue.activeAlertCount !== undefined) {
        count += issue.activeAlertCount;
      }
    }
    return count;
  }
}
