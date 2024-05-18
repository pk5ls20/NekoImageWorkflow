package sqlite

import (
	"errors"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	"github.com/sirupsen/logrus"
)

func InsertDbData[T dbDataModel](data *dbData[T]) error {
	if err := checkDBInstance(); err != nil {
		return commonLog.ErrorWrap(err)
	}
	if data == nil {
		return commonLog.ErrorWrap(errors.New("data which passed to InsertDbData is nil"))
	}
	var tmpData, err = encodeData(data)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	var tmpStruct = dbDataStoredModel{
		Tag:  data.Tag,
		Data: tmpData,
	}
	if err_ := dbInstance.Create(&tmpStruct); err_ != nil {
		return commonLog.ErrorWrap(err_.Error)
	}
	return nil
}

func InsertBatchDbData[T dbDataModel](data []*dbData[T]) error {
	if err := checkDBInstance(); err != nil {
		return commonLog.ErrorWrap(err)
	}
	if len(data) == 0 {
		logrus.Warning("*[]dbData which passed to InsertBatchDbData is empty")
		return nil
	}
	var tmpData = make([]*dbDataStoredModel, 0)
	for _, d := range data {
		encodedData, err := encodeData(d)
		if err != nil {
			return commonLog.ErrorWrap(err)
		}
		tmpData = append(tmpData, &dbDataStoredModel{
			Tag:  d.Tag,
			Data: encodedData,
		})
	}
	result := dbInstance.Create(&tmpData)
	if result.Error != nil {
		return commonLog.ErrorWrap(result.Error)
	}
	logrus.Debug("Successfully inserted ", result.RowsAffected, " records")
	return nil
}

func FindDbDataModelByTag(id keyTag) ([]*dbDataStoredModel, error) {
	var data = make([]*dbDataStoredModel, 0)
	if err := checkDBInstance(); err != nil {
		return data, commonLog.ErrorWrap(err)
	}
	if err := dbInstance.Where("Tag = ?", id).Find(&data); err.Error != nil {
		return data, commonLog.ErrorWrap(err.Error)
	}
	return data, nil
}

func DeleteDbDataByTag(id keyTag) error {
	if err := checkDBInstance(); err != nil {
		return commonLog.ErrorWrap(err)
	}
	err_ := dbInstance.Where("Tag = ?", id).Delete(&dbDataStoredModel{})
	if err_.Error != nil {
		return commonLog.ErrorWrap(err_.Error)
	}
	return nil
}

// checkDBInstance checks if the database is initialised
func checkDBInstance() error {
	if dbInstance == nil {
		return commonLog.ErrorWrap(errors.New("database not initialised"))
	}
	return nil
}
