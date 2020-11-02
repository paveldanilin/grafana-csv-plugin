import FileDatasource from './datasource';
import FileDatasourceQueryCtrl from './query_ctrl';

class FileConfigCtrl {
  constructor() {
    this.current.jsonData.encrypt = 'true';

    this.columnTypes = [
      { text: 'Text', value: 'text' },
      { text: 'Integer', value: 'integer' },
      { text: 'Real', value: 'real' },
      { text: 'Timestamp', value: 'timestamp' },
      { text: 'Date', value: 'date' },
    ];

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
    this.current.jsonData.columns = this.current.jsonData.columns || [];
  }

  addColumn(evt) {
    evt.preventDefault();
    this.current.jsonData.columns.push({
      name: '',
      type: 'text',
    });
  }

  deleteColumn(rowIndex) {
    this.current.jsonData.columns.splice(rowIndex, 1);
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
