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
          if (!this.current.jsonData.accessMode) {
            this.current.jsonData.accessMode = 'local';
          }

          if (!this.current.jsonData.filename) {
            this.current.jsonData.filename = '';
          }
          if (!this.current.jsonData.csvDelimiter) {
            this.current.jsonData.csvDelimiter = '';
          }
          if (!this.current.jsonData.csvComment) {
            this.current.jsonData.csvComment = '';
          }
          if (!this.current.jsonData.csvTrimLeadingSpace) {
            this.current.jsonData.csvTrimLeadingSpace = true;
          }

          if (!this.current.jsonData.sftpHost) {
            this.current.jsonData.sftpHost = '';
          }
          if (!this.current.jsonData.sftpPort) {
            this.current.jsonData.sftpPort = '';
          }
          if (!this.current.jsonData.sftpUser) {
            this.current.jsonData.sftpUser = '';
          }
          if (!this.current.jsonData.sftpIgnoreHostKey) {
            this.current.jsonData.sftpIgnoreHostKey = false;
          }
          if (!this.current.jsonData.sftpWorkingDir) {
            this.current.jsonData.sftpWorkingDir = '';
          }

          this.current.secureJsonData = this.current.secureJsonData || {};
          if (!this.current.secureJsonData.sftpPassword) {
            this.current.secureJsonData.sftpPassword = null;
          }

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
          key: 'onFilenameUpdate',
          value: function onFilenameUpdate() {
            this.current.url = this.current.jsonData.filename;
          }
        }, {
          key: 'onSftpHostUpdate',
          value: function onSftpHostUpdate() {
            this.current.url = this.current.jsonData.sftpHost + '>' + this.current.jsonData.filename;
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
