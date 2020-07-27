'use strict';

System.register(['./datasource', './query_ctrl'], function (_export, _context) {
  var FileDatasource, FileDatasourceQueryCtrl, _createClass, FileConfigCtrl, FileAnnotationsQueryCtrl;

  function _classCallCheck(instance, Constructor) {
    if (!(instance instanceof Constructor)) {
      throw new TypeError("Cannot call a class as a function");
    }
  }

  return {
    setters: [function (_datasource) {
      FileDatasource = _datasource.default;
    }, function (_query_ctrl) {
      FileDatasourceQueryCtrl = _query_ctrl.default;
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

      _export('ConfigCtrl', FileConfigCtrl = function () {
        function FileConfigCtrl() {
          var _this = this;

          _classCallCheck(this, FileConfigCtrl);

          this.current.jsonData.encrypt = 'true';

          this.accessModes = [{ text: 'Local', value: 'local' }, { text: 'SFTP', value: 'sftp' }];

          this.columnTypes = [{ text: 'Text', value: 'text' }, { text: 'Integer', value: 'integer' }, { text: 'Real', value: 'real' }, { text: 'Timestamp', value: 'timestamp' }, { text: 'Date', value: 'date' }];

          this.current.jsonData.accessMode = this.current.jsonData.accessMode || 'local';
          this.current.jsonData.filename = this.current.jsonData.filename || '';
          this.current.jsonData.csvDelimiter = this.current.jsonData.csvDelimiter || '';
          this.current.jsonData.csvComment = this.current.jsonData.csvComment || '';
          this.current.jsonData.csvTrimLeadingSpace = this.current.jsonData.csvTrimLeadingSpace || true;
          this.current.jsonData.sftpHost = this.current.jsonData.sftpHost || '';
          this.current.jsonData.sftpPort = this.current.jsonData.sftpPort || '';
          this.current.jsonData.sftpUser = this.current.jsonData.sftpUser || '';
          this.current.jsonData.sftpIgnoreHostKey = this.current.jsonData.sftpIgnoreHostKey || false;
          this.current.jsonData.sftpWorkingDir = this.current.jsonData.sftpWorkingDir || '';
          this.current.secureJsonData = this.current.secureJsonData || {};
          this.current.secureJsonData.sftpPassword = this.current.secureJsonData.sftpPassword || null;
          this.current.jsonData.columns = this.current.jsonData.columns || [];

          this.onPasswordReset = function (event) {
            event.preventDefault();
            _this.current['sftpPassword'] = null;
            _this.current.secureJsonFields['sftpPassword'] = false;
            _this.current.secureJsonData = _this.current.secureJsonData || {};
            _this.current.secureJsonData['sftpPassword'] = '';
          };

          this.onPasswordChange = function (event) {
            _this.current.secureJsonData = _this.current.secureJsonData || {};
            _this.current.secureJsonData['sftpPassword'] = event.currentTarget.value;
          };
        }

        _createClass(FileConfigCtrl, [{
          key: 'addColumn',
          value: function addColumn(evt) {
            evt.preventDefault();
            this.current.jsonData.columns.push({
              name: '',
              type: 'text'
            });
          }
        }, {
          key: 'deleteColumn',
          value: function deleteColumn(rowIndex) {
            this.current.jsonData.columns.splice(rowIndex, 1);
          }
        }]);

        return FileConfigCtrl;
      }());

      FileConfigCtrl.templateUrl = 'partials/config.html';

      _export('AnnotationsQueryCtrl', FileAnnotationsQueryCtrl = function FileAnnotationsQueryCtrl() {
        _classCallCheck(this, FileAnnotationsQueryCtrl);
      });

      _export('Datasource', FileDatasource);

      _export('QueryCtrl', FileDatasourceQueryCtrl);

      _export('ConfigCtrl', FileConfigCtrl);

      _export('AnnotationsQueryCtrl', FileAnnotationsQueryCtrl);
    }
  };
});
//# sourceMappingURL=module.js.map
