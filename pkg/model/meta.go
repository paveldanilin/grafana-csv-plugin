package model

type Meta struct {
	TableName string
	FileName string
	FileSize int64
	FileModTime int64
}

func NewMeta(tableName, fileName string, fileSize, fileModTime int64) *Meta {
	return &Meta{
		TableName: tableName,
		FileName: fileName,
		FileSize: fileSize,
		FileModTime: fileModTime,
	}
}
