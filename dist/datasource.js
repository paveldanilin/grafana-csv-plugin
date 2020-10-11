'use strict';

System.register(['lodash', './response_parser'], function (_export, _context) {
  var _, ResponseParse, _createClass, FileDatasource;

  function _classCallCheck(instance, Constructor) {
    if (!(instance instanceof Constructor)) {
      throw new TypeError("Cannot call a class as a function");
    }
  }

  return {
    setters: [function (_lodash) {
      _ = _lodash.default;
    }, function (_response_parser) {
      ResponseParse = _response_parser.default;
    }],
    execute: function () {
      _createClass = function () {
        function defineProperties(target, props) {
          for (var i = 0; i < props.length; i++) {
            var descriptor = props[i];
            descriptor.enumerable = descriptor.enumerable || false;
            descriptor.configurable = true;
            if ("value" in descriptor) descriptor.writable = true;
            Object.defineProperty(target, descriptor.key, descriptor);
          }
        }

        return function (Constructor, protoProps, staticProps) {
          if (protoProps) defineProperties(Constructor.prototype, protoProps);
          if (staticProps) defineProperties(Constructor, staticProps);
          return Constructor;
        };
      }();

      FileDatasource = function () {
        function FileDatasource(instanceSettings, backendSrv, timeSrv, templateSrv) {
          _classCallCheck(this, FileDatasource);

          this.id = instanceSettings.id;
          this.name = instanceSettings.name;
          this.backendSrv = backendSrv;
          this.timeSrv = timeSrv;
          this.templateSrv = templateSrv;
          this.responseParser = new ResponseParse();
        }

        _createClass(FileDatasource, [{
          key: 'query',
          value: function query(options) {
            var _this = this;

            // TODO: skip for filtering and process only first query
            var queries = _.filter(options.targets, function (target) {
              return target.hide !== true;
            }).map(function (target) {
              var rawQuery = target.query || _this.defaultSql();
              return {
                refId: target.refId,
                intervalMs: options.intervalMs,
                maxDataPoints: options.maxDataPoints,
                datasourceId: _this.id,
                format: target.format,
                query: _this.templateSrv.replace(rawQuery, _this.templateSrv.variables, _this.interpolateVar)
              };
            });

            if (queries.length === 0) {
              return Promise.resolve({
                data: []
              });
            }

            return this.backendSrv.datasourceRequest({
              url: '/api/tsdb/query',
              data: {
                from: options.range.from.valueOf().toString(),
                to: options.range.to.valueOf().toString(),
                // !!!!!!!!!!!!!!!!!!!!!!!!!
                // Perform only first query
                // !!!!!!!!!!!!!!!!!!!!!!!!!
                queries: [queries[0]]
              },
              method: 'POST'
            }).then(this.responseParser.processQueryResult);
          }
        }, {
          key: 'interpolateVar',
          value: function interpolateVar(value, variable) {
            if (typeof value === 'string') {
              if (variable.multi || variable.includeAll) {
                return value.replace(/'/g, '\'\'');
              } else {
                return value;
              }
            }

            if (typeof value === 'number') {
              return value;
            }

            var quotedValues = _.map(value, function (val) {
              if (typeof value === 'number') {
                return value;
              }
              return encodeURI(val.replace(/'/g, '\'\''));
            });

            return quotedValues.join(',');
          }
        }, {
          key: 'testDatasource',
          value: function testDatasource() {
            return this.backendSrv.datasourceRequest({
              url: '/api/tsdb/query',
              method: 'POST',
              data: {
                from: '5m',
                to: 'now',
                queries: [{
                  refId: '[tests-connection]',
                  intervalMs: 1,
                  maxDataPoints: 1,
                  datasourceId: this.id,
                  format: 'table',
                  query: ''
                }]
              }
            }).then(function (response) {
              if (response.status === 200) {
                return { status: 'success', message: 'Data source is working', title: 'Success' };
              }
            }).catch(function (err) {
              if (err.data && err.data.message) {
                return { status: 'error', message: err.data.message };
              } else {
                return { status: 'error', message: err.data.status };
              }
            });
          }
        }, {
          key: 'metricFindQuery',
          value: function metricFindQuery(query, optionalOptions) {
            var _this2 = this;

            var refId = 'tempvar';
            if (optionalOptions && optionalOptions.variable && optionalOptions.variable.name) {
              refId = optionalOptions.variable.name;
            }

            var interpolatedQuery = {
              refId: refId,
              datasourceId: this.id,
              query: this.templateSrv.replace(query, this.templateSrv.variables, this.interpolateVar),
              format: 'table'
            };

            var range = this.timeSrv.timeRange();
            var data = {
              queries: [interpolatedQuery],
              from: range.from.valueOf().toString(),
              to: range.to.valueOf().toString()
            };

            return this.backendSrv.datasourceRequest({
              url: '/api/tsdb/query',
              method: 'POST',
              data: data
            }).then(function (data) {
              return _this2.responseParser.parseMetricFindQueryResult(refId, data);
            });
          }
        }, {
          key: 'defaultSql',
          value: function defaultSql() {
            var defSql = 'SELECT * FROM {TableName} LIMIT 1, 15';
            return defSql.replace('{TableName}', this.name);
          }
        }]);

        return FileDatasource;
      }();

      _export('default', FileDatasource);
    }
  };
});
//# sourceMappingURL=datasource.js.map
