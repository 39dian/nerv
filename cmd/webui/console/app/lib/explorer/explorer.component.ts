import { Component, OnInit, ViewChild } from '@angular/core';
import { Router, ActivatedRoute, Params, UrlSegment } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { TreeComponent, TreeNode, ITreeOptions } from 'angular2-tree-component';
import { AceDirective } from '../form/ace/ace.directive';
import { FormModal } from '../form/form.modal';
import { ConfirmModal } from '../form/confirm.modal';
import { File, FileService } from './file.service';

const ADD_DIR_FORM = {
    name: "explorer_add_form",
    fields: [
        {
            name: "name", label: "名称", control: "text", type: "string", validators: { 'required': { message: '不能为空' } }
        },
        {
            name: "dirType", label: "类型", control: "select", type: "string",
            display: {
                options: [
                    { label: '子文件夹', value: 'child', default: true },
                    { label: '根文件夹', value: 'root' }
                ]
            }
        }
    ]
};

const ADD_FILE_FORM = {
    name: "explorer_add_form",
    fields: [
        {
            name: "name", label: "名称", control: "text", type: "string", validators: { 'required': { message: '不能为空' } }
        },
        {
            name: "fileType", label: "类型", control: "select", type: "string",
            display: {
                options: [
                    { label: 'JSON', value: 'json', default: true },
                    { label: 'Shell', value: 'sh' },
                    { label: 'Ruby', value: 'ruby' },
                    { label: 'Python', value: 'python' }
                ]
            }
        }
    ]
};

const RENAME_FORM = {
    name: "explorer_rename_form",
    fields: [
        {
            name: "name", label: "名称", control: "text", type: "string", validators: { 'required': { message: '不能为空' } }
        }
    ]
};

@Component({
    selector: 'nerv-explorer',
    templateUrl: 'app/lib/explorer/explorer.component.html',
})
export class ExplorerComponent implements OnInit {
    options = {
        idField: 'name',
        displayField: 'title',
        getChildren: (node: TreeNode) => {
            let name = node.data.name;
            let parent = node.parent;
            while (parent && parent.data && parent.data['type']) {
                name = parent.data.name + "/" + name;
                parent = parent.parent;
            }
            return this.fileService.get(name);
        }
    };

    nodes = [];

    @ViewChild(TreeComponent)
    private tree: TreeComponent;

    @ViewChild(AceDirective)
    private ace: AceDirective;

    selectedNode: File

    get canRemove(): boolean {
        return this.selectedNode != null;
    }

    get canRename(): boolean {
        return this.selectedNode != null;
    }

    get canAddFile(): boolean {
        return this.selectedNode != null;
    }

    get canSave(): boolean {
        return this.selectedNode && this.selectedNode.type == 'file';
    }

    get mode(): string {
        return this.selectedNode ? this.selectedNode.extension : 'json';
    }

    constructor(
        private modalService: NgbModal,
        private fileService: FileService
    ) {
    }

    ngOnInit(): void {


        this.fileService.get('')
            .then((data) => {
                this.nodes = data;
            })
            .catch((error) => this.error('加载错误', `加载目录失败\r\n${error}`));
    }

    onAddFile() {
        let newItem = {};
        const modalRef = this.modalService.open(FormModal, { backdrop: 'static' });
        modalRef.componentInstance.title = '添加文件';
        modalRef.componentInstance.form = ADD_FILE_FORM;
        modalRef.componentInstance.data = newItem;
        modalRef.result.then((result) => {
            if (result == 'ok') {
                this.addFile(this.selectedNode, newItem);
            }
        });
    }

    onAddDir() {
        let newItem = {};
        const modalRef = this.modalService.open(FormModal, { backdrop: 'static' });
        modalRef.componentInstance.title = '添加文件夹';
        modalRef.componentInstance.form = ADD_DIR_FORM;
        modalRef.componentInstance.data = newItem;
        modalRef.result.then((result) => {
            if (result == 'ok') {
                this.addDir(this.selectedNode, newItem);
            }
        });
    }

    onRemove() {
        if (!this.canRemove) return;
        const modalRef = this.modalService.open(ConfirmModal, { backdrop: 'static' });
        modalRef.componentInstance.title = '删除';
        modalRef.componentInstance.message = `删除${this.selectedNode.name}?`;
        modalRef.result.then((result) => {
            if (result == 'ok') {
                this.remove(this.selectedNode);
            }
        });
    }



    onRename() {
        let newItem = {};
        const modalRef = this.modalService.open(FormModal, { backdrop: 'static' });
        modalRef.componentInstance.title = '重命名';
        modalRef.componentInstance.form = RENAME_FORM;
        modalRef.componentInstance.data = newItem;
        modalRef.result.then((result) => {
            if (result == 'ok') {
                this.rename(newItem['name']);
            }
        });
    }

    onSave() {
        const file = this.selectedNode;
        if (file && file.type == 'file') {
            this.saveFile(file);
        }
    }

    onTextChanged(text: string) {
        const file = this.selectedNode;
        if (file && file.type == 'file') {
            file.dirty = true;
            file.content = this.getAceContent();
            this.tree.treeModel.update();
        }
    }

    onActive(event: {}) {
        const node = event['node'];
        const file = this.selectedNode = node['data'];
        if (file && file.type == 'file') {
            this.loadContent(file);
        }
    }

    private saveFile(file: File): void {

    }

    private getAceContent(): string {
        return this.ace.text;
    }

    private loadContent(file: File): void {
        const content = file.content;
        this.ace.text = content || '';
    }

    private rename(name: string): void {
        this.selectedNode.name = name;
        this.tree.treeModel.update();
    }

    private remove(file: File): void {
        const activeNode = this.tree.treeModel.activeNodes[0];
        const index = activeNode.parent.data.children.indexOf(activeNode.data);
        activeNode.parent.data.children.splice(index, 1);
        this.tree.treeModel.update();
        this.ace.text = '';
    }

    private addFile(parent: File, node: {}): void {
        const selectedNode = this.selectedNode;
        const name = node['name'];
        if (selectedNode) {
            if (selectedNode.type == 'dir') {
                let children = selectedNode.children;
                if (!children) {
                    this.selectedNode.children = children = new Array<File>();
                }

                children.push({ name: name, type: 'file', extension: node['fileType'] });
                this.tree.treeModel.update();
            } else {
                let children = this.tree.treeModel.activeNodes[0].parent.data.children;
                const title = `${name}.${node['fileType']}`;
                children.push({ name: title, type: 'file', title: title, extension: node['fileType'] });
                this.tree.treeModel.update();
            }
        }
    }

    private addDir(parent: File, node: {}): void {
        const name = node['name'];
        const selectedNode = this.selectedNode;
        let nodes;
        if (selectedNode) {
            if (selectedNode['type'] == 'file') {
                nodes = this.tree.treeModel.activeNodes[0].parent.data.children;
            } else if (node['dirType'] == 'child') {
                nodes = selectedNode.children;
                if (!nodes) {
                    selectedNode.children = nodes = new Array<File>();
                }
            } else {
                nodes = this.nodes;
            }
        } else {
            nodes = this.nodes;
        }
        nodes.push({ name: name, type: 'dir' });
        this.tree.treeModel.update();
    }

    private error(title: string, error: any): void {
        const modalRef = this.modalService.open(ConfirmModal, { backdrop: 'static' });
        modalRef.componentInstance.title = title;
        modalRef.componentInstance.message = error.toString();
        modalRef.componentInstance.buttons = { ok: true, cancel: false };
    }
}