import _ from 'lodash';
import ResponseParse from './response_parser';

export default class FileDatasource {
  constructor(instanceSettings, backendSrv, timeSrv, templateSrv, variableSrv) {
    this.id = instanceSettings.id;
    this.name = instanceSettings.name;
    this.backendSrv = backendSrv;
    this.timeSrv = timeSrv;
    this.templateSrv = templateSrv;
    this.variableSrv = variableSrv;
    this.responseParser = new ResponseParse();
  }

  query(options) {
    // TODO: skip for filtering and process only first query
    const queries = _.filter(options.targets, target => {
      return target.hide !== true;
    }).map(target => {
      return {
        refId: target.refId,
        intervalMs: options.intervalMs,
        maxDataPoints: options.maxDataPoints,
        datasourceId: this.id,
        format: target.format,
        // query: target.query || this.defaultSql(),
        query: this.templateSrv.replace(target.query, this.variableSrv.variables, this.interpolateVar),
      };
    });

    if (queries.length === 0) {
      return Promise.resolve({
        data: [],
      });
    }

    return this.backendSrv.datasourceRequest({
      url: `/api/tsdb/query`,
      data: {
        from: options.range.from.valueOf().toString(),
        to: options.range.to.valueOf().toString(),
        // !!!!!!!!!!!!!!!!!!!!!!!!!
        // Perform only first query
        // !!!!!!!!!!!!!!!!!!!!!!!!!
        queries: [queries[0]],
      },
      method: 'POST',
    }).then(this.responseParser.processQueryResult);
  }

  interpolateVar(value, variable) {
    if (typeof value === 'string') {
      if (variable.multi || variable.includeAll) {
        return value.replace(/'/g, `''`);
      } else {
        return value;
      }
    }

    if (typeof value === 'number') {
      return value;
    }

    const quotedValues = _.map(value, val => {
      if (typeof value === 'number') {
        return value;
      }
      return encodeURI(val.replace(/'/g, `''`));
    });

    return quotedValues.join(',');
  }

  testDatasource() {
    return this.backendSrv.datasourceRequest({
      url: '/api/tsdb/query',
      method: 'POST',
      data: {
        from: '5m',
        to: 'now',
        queries: [
          {
            refId: '[tests-connection]',
            intervalMs: 1,
            maxDataPoints: 1,
            datasourceId: this.id,
            format: 'table',
            query: '',
          }
        ]
      }
    }).then((response) => {
      if (response.status === 200) {
        return { status: 'success', message: 'Data source is working', title: 'Success' };
      }
    }).catch((err) => {
      if (err.data && err.data.message) {
        return { status: 'error', message: err.data.message };
      } else {
        return { status: 'error', message: err.data.status };
      }
    });
  }

  metricFindQuery(query, optionalOptions) {
    let refId = 'mqtmp';
    if (optionalOptions && optionalOptions.variable && optionalOptions.variable.name) {
      refId = optionalOptions.variable.name;
    }
    const interpolatedQuery = {
      refId: refId,
      datasourceId: this.id,
      format: 'table',
      query: '',
    };

    const range = this.timeSrv.timeRange();
    const data = {
      queries: [interpolatedQuery],
      from: range.from.valueOf().toString(),
      to: range.to.valueOf().toString(),
    };

    return this.backendSrv.datasourceRequest({
        url: '/api/tsdb/query',
        method: 'POST',
        data: data,
      })
      .then((data) => this.responseParser.parseMetricFindQueryResult(refId, data));
  }

  defaultSql() {
    const defSql = 'SELECT * FROM {TableName} LIMIT 1, 15';
    return defSql.replace('{TableName}', this.name);
  }
}
