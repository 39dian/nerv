import { Component, ViewContainerRef, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { FormConfig } from '../config/form.config';
import { FormRegistry } from '../form/form.registry';
import { RestyService } from '../resty/resty.service';
import { FormBaseComponent } from './base.form';
import { Form, Field } from '../form/model';


@Component({
    templateUrl: 'app/lib/app/detail.form.html',
})
export class DetailComponent extends FormBaseComponent {

    constructor(
        formRegistry: FormRegistry,
        configService: FormConfig,
        modalService: NgbModal,
        router: Router,
        route: ActivatedRoute,
        resty: RestyService
    ) {
        super('detail', formRegistry, configService, modalService, router, route, resty);
    }
}