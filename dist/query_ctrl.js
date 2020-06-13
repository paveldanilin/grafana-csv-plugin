'use strict';

System.register(['lodash', 'app/plugins/sdk', '@grafana/data'], function (_export, _context) {
  var _, QueryCtrl, PanelEvents, _createClass, FileDatasourceQueryCtrl;

  function _classCallCheck(instance, Constructor) {
    if (!(instance instanceof Constructor)) {
      throw new TypeError("Cannot call a class as a function");
    }
  }

  function _possibleConstructorReturn(self, call) {
    if (!self) {
      throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
    }

    return call && (typeof call === "object" || typeof call === "function") ? call : self;
  }

  function _inherits(subClass, superClass) {
    if (typeof superClass !== "function" && superClass !== null) {
      throw new TypeError("Super expression must either be null or a function, not " + typeof superClass);
    }

    subClass.prototype = Object.create(superClass && superClass.prototype, {
      constructor: {
        value: subClass,
        enumerable: false,
        writable: true,
        configurable: true
      }
    });
    if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass;
  }

  return {
    setters: [function (_lodash) {
      _ = _lodash.default;
    }, function (_appPluginsSdk) {
      QueryCtrl = _appPluginsSdk.QueryCtrl;
    }, function (_grafanaData) {
      PanelEvents = _grafanaData.PanelEvents;
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

      FileDatasourceQueryCtrl = function (_QueryCtrl) {
        _inherits(FileDatasourceQueryCtrl, _QueryCtrl);

        function FileDatasourceQueryCtrl($scope, $injector) {
          _classCallCheck(this, FileDatasourceQueryCtrl);

          var _this = _possibleConstructorReturn(this, (FileDatasourceQueryCtrl.__proto__ || Object.getPrototypeOf(FileDatasourceQueryCtrl)).call(this, $scope, $injector));

          _this.target.alias = '';
          _this.target.query = _this.target.query || 'SELECT * FROM DataTable LIMIT 1, 15';

          // FORMAT
          _this.target.format = _this.target.format || 'table';
          _this.formats = [{ text: 'Time series', value: 'time_series' }, { text: 'Table', value: 'table' }];

          _this.panelCtrl.events.on(PanelEvents.dataReceived, _this.onDataReceived.bind(_this), $scope);
          _this.panelCtrl.events.on(PanelEvents.dataError, _this.onDataError.bind(_this), $scope);
          return _this;
        }

        _createClass(FileDatasourceQueryCtrl, [{
          key: 'onDataReceived',
          value: function onDataReceived(dataList) {
            this.lastQueryMeta = null;
            this.lastQueryError = null;
            var anySeriesFromQuery = _.find(dataList, { refId: this.target.refId });
            if (anySeriesFromQuery) {
              this.lastQueryMeta = anySeriesFromQuery.meta;
            }
          }
        }, {
          key: 'onDataError',
          value: function onDataError(err) {
            if (err.data && err.data.results) {
              var queryRes = err.data.results[this.target.refId];
              if (queryRes) {
                this.lastQueryMeta = queryRes.meta;
                this.lastQueryError = queryRes.error;
              }
            }
          }
        }]);

        return FileDatasourceQueryCtrl;
      }(QueryCtrl);

      _export('default', FileDatasourceQueryCtrl);

      FileDatasourceQueryCtrl.templateUrl = 'partials/query.editor.html';
    }
  };
});
//# sourceMappingURL=query_ctrl.js.map
