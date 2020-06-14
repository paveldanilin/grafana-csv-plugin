package csv

import (
	"github.com/hashicorp/go-hclog"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"sync"
)

type metaDB struct {
	db *DB
	descriptor *FileDescriptor
}

type DbManager interface {
	Get(dbName string, descriptor *FileDescriptor) (*DB, error)
}

type dbManagerImpl struct {
	mux sync.Mutex
	databases map[string]*metaDB
	logger hclog.Logger
}

func NewDbManager(logger hclog.Logger) DbManager {
	return &dbManagerImpl{
		mux:       sync.Mutex{},
		databases: make(map[string]*metaDB),
		logger: logger,
	}
}

func (m *dbManagerImpl) Get(dbName string, descriptor *FileDescriptor) (*DB, error) {
	m.logger.Debug("Going to obtain DB", "dbName", dbName, "filename", descriptor.Filename)
	m.mux.Lock()
	defer m.mux.Unlock()

	if !m.has(dbName) {
		m.logger.Debug("Loading CSV file", "dbName", dbName, "filename", descriptor.Filename)
		db, err := Read(descriptor)
		if err != nil {
			m.logger.Error("Could not load CSV file", "error", err.Error(), "dbName", dbName, "filename", descriptor.Filename)
			return nil, err
		}
		m.databases[dbName] = &metaDB{
			db:         db,
			descriptor: descriptor,
		}
		m.logger.Debug("CSV has been loaded", "dbName", dbName, "filename", descriptor.Filename)
	} else {
		m.logger.Debug("CSV has been loaded, going to detect file change", "dbName", dbName, "filename", descriptor.Filename)
		filename := m.databases[dbName].descriptor.Filename
		fSize := m.databases[dbName].descriptor.fileSize
		fModTime := m.databases[dbName].descriptor.fileModTime
		isCsvChanged, err := util.FileChanged(filename, fSize, fModTime)
		if err != nil {
			m.logger.Error("Could not detect file change", "error", err.Error(), "filename", filename)
			return nil, err
		}
		if isCsvChanged {
			m.logger.Debug("CSV file has been changed, going to reload data", "dbName", dbName, "filename", descriptor.Filename)
			m.databases[dbName].db.Release()
			m.databases[dbName].db = nil

			db, err := Read(m.databases[dbName].descriptor)
			if err != nil {
				m.logger.Error("Could not reload CSV file", "dbName", dbName, "error", err.Error(), "filename", descriptor.Filename)
				return nil, err
			}
			m.databases[dbName].db = db
		} else {
			m.logger.Debug("CSV has not been changed, going to use in-memory DB", "dbName", dbName, "filename", descriptor.Filename)
		}
	}

	m.logger.Info("CSV has been loaded into memory", "dbName", dbName, "filename", descriptor.Filename)
	return m.databases[dbName].db, nil
}

func (m *dbManagerImpl) has(dbName string) bool {
	_, ok := m.databases[dbName]
	return ok && m.databases[dbName].db != nil
}
