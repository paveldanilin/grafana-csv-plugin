import FileDatasource from './datasource';
import FileDatasourceQueryCtrl from './query_ctrl';

class FileConfigCtrl {
  constructor() {
    this.current.jsonData.encrypt = 'true';

    this.accessModes = [
      { text: 'Local', value: 'local' },
      { text: 'SFTP', value: 'sftp' },
    ];
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

    this.onPasswordReset = (event) => {
      event.preventDefault();
      this.current['sftpPassword'] = null;
      this.current.secureJsonFields['sftpPassword'] = false;
      this.current.secureJsonData = this.current.secureJsonData || {};
      this.current.secureJsonData['sftpPassword'] = '';
    };

    this.onPasswordChange = (event) => {
      this.current.secureJsonData =  this.current.secureJsonData || {};
      this.current.secureJsonData['sftpPassword'] = event.currentTarget.value;
    };
  }

  onFilenameUpdate() {
    this.current.url = this.current.jsonData.filename;
  }

  onSftpHostUpdate() {
    this.current.url = this.current.jsonData.sftpHost + '>' + this.current.jsonData.filename;
  }
}
FileConfigCtrl.templateUrl = 'partials/config.html';

class FileAnnotationsQueryCtrl {
  constructor() {
  }
}

export {
  FileDatasource as Datasource,
  FileDatasourceQueryCtrl as QueryCtrl,
  FileConfigCtrl as ConfigCtrl,
  FileAnnotationsQueryCtrl as AnnotationsQueryCtrl,
};
