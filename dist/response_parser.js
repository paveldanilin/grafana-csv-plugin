'use strict';

System.register(['lodash'], function (_export, _context) {
  var _, _createClass, ResponseParse;

  function _classCallCheck(instance, Constructor) {
    if (!(instance instanceof Constructor)) {
      throw new TypeError("Cannot call a class as a function");
    }
  }

  return {
    setters: [function (_lodash) {
      _ = _lodash.default;
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

      ResponseParse = function () {
        function ResponseParse() {
          _classCallCheck(this, ResponseParse);
        }

        _createClass(ResponseParse, [{
          key: 'processQueryResult',
          value: function processQueryResult(res) {
            var data = [];

            if (!res.data.results) {
              return { data: data };
            }

            for (var key in res.data.results) {
              var queryRes = res.data.results[key];

              if (queryRes.series) {
                var _iteratorNormalCompletion = true;
                var _didIteratorError = false;
                var _iteratorError = undefined;

                try {
                  for (var _iterator = queryRes.series[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true) {
                    var series = _step.value;

                    data.push({
                      target: series.name,
                      datapoints: series.points,
                      refId: queryRes.refId,
                      meta: queryRes.meta
                    });
                  }
                } catch (err) {
                  _didIteratorError = true;
                  _iteratorError = err;
                } finally {
                  try {
                    if (!_iteratorNormalCompletion && _iterator.return) {
                      _iterator.return();
                    }
                  } finally {
                    if (_didIteratorError) {
                      throw _iteratorError;
                    }
                  }
                }
              }

              if (queryRes.tables) {
                var _iteratorNormalCompletion2 = true;
                var _didIteratorError2 = false;
                var _iteratorError2 = undefined;

                try {
                  for (var _iterator2 = queryRes.tables[Symbol.iterator](), _step2; !(_iteratorNormalCompletion2 = (_step2 = _iterator2.next()).done); _iteratorNormalCompletion2 = true) {
                    var table = _step2.value;

                    table.type = 'table';
                    table.refId = queryRes.refId;
                    table.meta = queryRes.meta;
                    data.push(table);
                  }
                } catch (err) {
                  _didIteratorError2 = true;
                  _iteratorError2 = err;
                } finally {
                  try {
                    if (!_iteratorNormalCompletion2 && _iterator2.return) {
                      _iterator2.return();
                    }
                  } finally {
                    if (_didIteratorError2) {
                      throw _iteratorError2;
                    }
                  }
                }
              }
            }

            return { data: data };
          }
        }, {
          key: 'parseMetricFindQueryResult',
          value: function parseMetricFindQueryResult(refId, results) {
            if (!results || results.data.length === 0 || results.data.results[refId].meta.rowCount === 0) {
              return [];
            }

            var columns = results.data.results[refId].tables[0].columns;
            var rows = results.data.results[refId].tables[0].rows;
            var textColIndex = this.findColIndex(columns, '__text');
            var valueColIndex = this.findColIndex(columns, '__value');

            if (columns.length === 2 && textColIndex !== -1 && valueColIndex !== -1) {
              return this.transformToKeyValueList(rows, textColIndex, valueColIndex);
            }

            return this.transformToSimpleList(rows);
          }
        }, {
          key: 'transformToKeyValueList',
          value: function transformToKeyValueList(rows, textColIndex, valueColIndex) {
            var res = [];

            for (var i = 0; i < rows.length; i++) {
              if (!this.containsKey(res, rows[i][textColIndex])) {
                res.push({
                  text: rows[i][textColIndex],
                  value: rows[i][valueColIndex]
                });
              }
            }

            return res;
          }
        }, {
          key: 'transformToSimpleList',
          value: function transformToSimpleList(rows) {
            var res = [];

            for (var i = 0; i < rows.length; i++) {
              for (var j = 0; j < rows[i].length; j++) {
                var value = rows[i][j];
                if (res.indexOf(value) === -1) {
                  res.push(value);
                }
              }
            }

            return _.map(res, function (value) {
              return { text: value };
            });
          }
        }, {
          key: 'findColIndex',
          value: function findColIndex(columns, colName) {
            for (var i = 0; i < columns.length; i++) {
              if (columns[i].text === colName) {
                return i;
              }
            }

            return -1;
          }
        }, {
          key: 'containsKey',
          value: function containsKey(res, key) {
            for (var i = 0; i < res.length; i++) {
              if (res[i].text === key) {
                return true;
              }
            }
            return false;
          }
        }, {
          key: 'transformAnnotationResponse',
          value: function transformAnnotationResponse(options, data) {
            var table = data.data.results[options.annotation.name].tables[0];

            var timeColumnIndex = -1;
            var timeEndColumnIndex = -1;
            var textColumnIndex = -1;
            var tagsColumnIndex = -1;

            for (var i = 0; i < table.columns.length; i++) {
              if (table.columns[i].text === 'time_sec' || table.columns[i].text === 'time') {
                timeColumnIndex = i;
              } else if (table.columns[i].text === 'timeend') {
                timeEndColumnIndex = i;
              } else if (table.columns[i].text === 'title') {
                return Promise.reject({
                  message: 'The title column for annotations is deprecated, now only a column named text is returned'
                });
              } else if (table.columns[i].text === 'text') {
                textColumnIndex = i;
              } else if (table.columns[i].text === 'tags') {
                tagsColumnIndex = i;
              }
            }

            if (timeColumnIndex === -1) {
              return Promise.reject({
                message: 'Missing mandatory time column (with time_sec column alias) in annotation query.'
              });
            }

            var list = [];
            for (var _i = 0; _i < table.rows.length; _i++) {
              var row = table.rows[_i];
              var timeEnd = timeEndColumnIndex !== -1 && row[timeEndColumnIndex] ? Math.floor(row[timeEndColumnIndex]) : undefined;
              list.push({
                annotation: options.annotation,
                time: Math.floor(row[timeColumnIndex]),
                timeEnd: timeEnd,
                text: row[textColumnIndex] ? row[textColumnIndex].toString() : '',
                tags: row[tagsColumnIndex] ? row[tagsColumnIndex].trim().split(/\s*,\s*/) : []
              });
            }

            return list;
          }
        }]);

        return ResponseParse;
      }();

      _export('default', ResponseParse);
    }
  };
});
//# sourceMappingURL=response_parser.js.map
