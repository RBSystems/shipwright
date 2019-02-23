import { Component, OnInit, ViewChild, ComponentFactoryResolver, forwardRef, Input } from '@angular/core';

import { DashPanelDirective } from './dashpanel.directive';
import { DashPanelService } from 'src/app/services/dashpanel.service';
import { IDashPanel } from './idashpanel';
import { DataService } from 'src/app/services/data.service';
import { AllAlerts, CritAlerts, WarnAlerts } from 'src/app/objects/alerts';

@Component({
  selector: 'dashpanel',
  templateUrl: './dashpanel.component.html',
  styleUrls: ['./dashpanel.component.scss']
})

export class DashPanelComponent implements OnInit {
  @Input() panelType: string = AllAlerts;

  @ViewChild(forwardRef(()=>DashPanelDirective)) direct: DashPanelDirective;

  constructor(private resolver: ComponentFactoryResolver, private dashServe: DashPanelService, public data: DataService) {
    
  }

  ngOnInit() {
    if(this.data.finished) {
      this.CreatePanel();
    } else {
      this.data.loaded.subscribe(() => {
        this.CreatePanel();
      })
    }
  }

  CreatePanel(info?: any) {
    let panel = this.dashServe.getPanel(this.panelType);
    let factory = this.resolver.resolveComponentFactory(panel.component);
    let viewContainerRef = this.direct.viewContainerRef;
    viewContainerRef.clear();

    let componentRef = viewContainerRef.createComponent(factory);

    if(info != null) {
      (<IDashPanel>componentRef.instance).info = info;
    }
  }
}
