import _ from 'lodash';
import { QueryCtrl } from 'app/plugins/sdk';
import { PanelEvents } from '@grafana/data';

export default class FileDatasourceQueryCtrl extends QueryCtrl {
  constructor($scope, $injector) {
    super($scope, $injector);

    this.tableColumns = []; // hold a full list of columns from server response
    this.searchableTableColumns = [];
    this.addColumnName = '';
    this.target.alias = '';

    // FORMAT
    this.target.format = this.target.format || 'table';
    this.formats = [
      { text: 'Time series', value: 'time_series' },
      { text: 'Table', value: 'table' },
    ];

    // COLUMN SORT
    this.target.columnSort = this.target.columnSort || 'ASC';
    this.target.customColumnOrder = this.target.customColumnOrder || [];
    this.columnSorts = [
      { text: 'Asc', value: 'ASC' },
      { text: 'Desc', value: 'DESC' },
      { text: 'Custom', value: 'CUSTOM' },
    ];

    this.panelCtrl.events.on(PanelEvents.dataReceived, this.onDataReceived.bind(this), $scope);
    this.panelCtrl.events.on(PanelEvents.dataError, this.onDataError.bind(this), $scope);
  }

  addColumn() {
    if (this.addColumnName.trim().length === 0) {
      return;
    }
    // Add to query
    this.target.customColumnOrder.push(this.addColumnName);
    // Remove from searchable
    const idx = this.searchableTableColumns.indexOf(this.addColumnName);
    if (idx !== -1) {
      this.searchableTableColumns.splice(idx, 1);
    }
    // Drop value
    this.addColumnName = '';
    this.refresh();
  }

  removeColumn(columnText) {
    if (columnText.trim().length === 0) {
      return;
    }
    // Remove from searchable
    this.searchableTableColumns.push(columnText);
    // Remove from query
    const idx = this.target.customColumnOrder.indexOf(columnText);
    if (idx !== -1) {
      this.target.customColumnOrder.splice(idx, 1);
    }
    this.refresh();
  }

  onDataReceived(dataList) {
    this.lastQueryMeta = null;
    this.lastQueryError = null;

    const anySeriesFromQuery = _.find(dataList, { refId: this.target.refId });
    if (anySeriesFromQuery) {
      this.lastQueryMeta = anySeriesFromQuery.meta;
    }

    this.tableColumns = [];
    const tableMeta = dataList.filter((data) => {
      return data.refId === this.target.refId;
    })[0] || null;

    if (tableMeta) {
      this.tableColumns = tableMeta.columns.map((col) => {
        return col.text;
      });
    }

    this.searchableTableColumns = this.tableColumns.filter((col) => {
      return this.target.customColumnOrder.indexOf(col) === -1;
    });
  }

  onDataError(err) {
    if (err.data && err.data.results) {
      const queryRes = err.data.results[this.target.refId];
      if (queryRes) {
        this.lastQueryMeta = queryRes.meta;
        this.lastQueryError = queryRes.error;
      }
    }
  }
}

FileDatasourceQueryCtrl.templateUrl = 'partials/query.editor.html';
