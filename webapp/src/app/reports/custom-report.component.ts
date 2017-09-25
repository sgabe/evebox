/* Copyright (c) 2017 Jason Ish
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

import {Component, ViewChild} from '@angular/core';
import {ApiService} from '../api.service';

@Component({
    template: `
      <div class="content">

        <!-- Dummy form that is the same height as the form below, so when it flashes, the page doesn't adjust in height. -->
        <form *ngIf="!formVisible" class="form-inline">
          <input class="form-control">
        </form>

        <form *ngIf="formVisible" class="form-inline" #inputForm="ngForm"
              (ngSubmit)="refresh()">

          <select class="form-control" [(ngModel)]="inputModel.reportType"
                  name="reportType">
            <option value="customAggregation">Custom Aggregation</option>
            <option value="presetAggregation">Preset Aggregation</option>
          </select>

          <div style="display:inline;"
               *ngIf="inputModel.reportType == 'customAggregation'">

            <div class="form-group">
              <label for="inputModelField">Field:</label>
              <input id="inputModelField" class="form-control"
                     [(ngModel)]="inputModel.field"
                     placeholder="Field..." name="field">
              <div class="dropdown" style="display:inline-block;">
                <button type="button"
                        class="btn btn-default dropdown-toggle"
                        data-toggle="dropdown"><span
                    class="caret"></span></button>
                <ul class="dropdown-menu">
                  <li><a href="javascript:void(0);"
                         (click)="setField($event)">tls.sni</a></li>
                  <li><a href="javascript:void(0);"
                         (click)="setField($event)">tls.subject</a>
                  </li>
                </ul>
              </div>
            </div>

            <div class="form-group">
              <label for="inputModel-size">Size:</label>
              <input id="inputModel-size" class="form-control"
                     [(ngModel)]="inputModel.size"
                     placeholder="Size..." name="size">
            </div>

            <div class="form-group">
              <label for="inputModel-order">Order:</label>
              <select class="form-control" [(ngModel)]="inputModel.order"
                      name="inputModel-order">
                <option value="DESC">DESC</option>
                <option value="ASC">ASC</option>
              </select>
            </div>

          </div>


          <button class="form-control btn"
                  [ngClass]="{'btn-primary': inputForm.dirty, 'btn-default': !inputForm.dirty}"
                  type="submit">
            Submit
          </button>

        </form>

        <br/>

        <div *ngIf="resultModel">
          <report-data-table
              [title]="resultModel.title"
              [rows]="resultModel.rows"
              [headers]="['#', resultModel.label]"></report-data-table>
        </div>

      </div>`
})
export class CustomReportComponent {

    @ViewChild("inputForm") inputForm: any;

    public formVisible: boolean = true;

    public inputModel: any = {
        reportType: "customAggregation",
        field: "",
        size: 10,
        order: "DESC",
    };

    public resultModel: any = null;

    results: any[] = null;

    constructor(private api: ApiService) {
    }

    setField(event: any) {
        this.inputModel.field = event.srcElement.innerText;
        this.inputForm.control.markAsDirty();
    }

    refresh() {
        this.api.reportAgg(this.inputModel.field, {
            size: this.inputModel.size,
            order: this.inputModel.order,
        }).then((response: any) => {
            this.resultModel = {
                rows: response.data,
                title: this.inputModel.field,
                label: this.inputModel.field,
            };
        });

        // Silly way of setting form non-dirty.
        this.formVisible = false;
        setTimeout(() => {
            this.formVisible = true;
        });
    }
}