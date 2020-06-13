package csv

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"sync"
)

type metaDB struct {
	db *DB
	descriptor *FileDescriptor
}

type DbManager struct {
	mux sync.Mutex
	databases map[string]*metaDB
}

func NewDbManager() *DbManager {
	return &DbManager{
		mux:       sync.Mutex{},
		databases: make(map[string]*metaDB),
	}
}

func (m *DbManager) Get(dbName string, descriptor *FileDescriptor) (*DB, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if !m.has(dbName) {
		db, err := Read(descriptor)
		if err != nil {
			return nil, err
		}
		m.databases[dbName] = &metaDB{
			db:         db,
			descriptor: descriptor,
		}
	} else {
		filename := m.databases[dbName].descriptor.Filename
		fSize := m.databases[dbName].descriptor.fileSize
		fModTime := m.databases[dbName].descriptor.fileModTime
		isCsvChanged, err := util.FileChanged(filename, fSize, fModTime)
		if err != nil {
			return nil, err
		}
		if isCsvChanged {
			m.databases[dbName].db.Release()
			m.databases[dbName].db = nil

			db, err := Read(m.databases[dbName].descriptor)
			if err != nil {
				return nil, err
			}
			m.databases[dbName].db = db
		}
	}
	return m.databases[dbName].db, nil
}

func (m *DbManager) has(dbName string) bool {
	_, ok := m.databases[dbName]
	return ok && m.databases[dbName].db != nil
}
