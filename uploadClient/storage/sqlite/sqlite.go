package sqlite

import (
	"fmt"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/bridge"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

type keyTag int

var (
	dbInstance *gorm.DB
	initDBOnce sync.Once
)

var dbName = "client_data.db"

// InitSqlite TODO: impl tx
// InitSqlite Initialises the database and writes the residual buffered data within the database to fileTransBridge
func InitSqlite() {
	initDBOnce.Do(func() {
		// 1. load (and create if not exists) sqlite database
		exe, err := os.Executable()
		if err != nil {
			logrus.Fatal(err)
		}
		exePath := filepath.Dir(exe)
		dbPath := filepath.Join(exePath, dbName)
		db, _err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if _err != nil {
			logrus.Fatal(_err)
		}
		dbInstance = db
		if err = dbInstance.AutoMigrate(&dbDataStoredModel{}); err != nil {
			logrus.Fatal(err)
		}
		// 2. load data into builtin channel
		receiveData, err := FindDbDataModelByTag(BridgeDataTag)
		if err != nil {
			logrus.Fatal(err)
		}
		tmpBytes := make([][]byte, len(receiveData))
		for i, data := range receiveData {
			tmpBytes[i] = data.Data
		}
		decodedReceiveData, err := decodeDataBatch[model.BridgeDataModel](tmpBytes)
		if err != nil {
			logrus.Fatal(err)
		}
		receiveDataToBridge[model.BridgeDataModel](decodedReceiveData)
		// 3. delete data from sqlite
		if err_ := DeleteDbDataByTag(BridgeDataTag); err_ != nil {
			logrus.Fatal(err_)
		}
	})
}

// CloseSqlite TODO: impl tx
// CloseSqlite Securely close the database and write the contents of fileTransBridge
func CloseSqlite() {
	if dbInstance != nil {
		// TODO: maybe we can move this out of sqlite.go, Consider dynamic registration functions
		// 1. write builtin channel data to database
		preUploadBridgeChannel := bridge.GetPreUploadTransBridgeInstance()
		if preUploadBridgeChannel != nil {
			tmpData, _ := preUploadBridgeChannel.PopAll()
			tmpDBData := make([]*dbData[model.BridgeDataModel], len(tmpData))
			for i, data := range tmpData {
				tmpDBData[i] = &dbData[model.BridgeDataModel]{Tag: BridgeDataTag, Data: data}
			}
			if err := InsertBatchDbData(tmpDBData); err != nil {
				logrus.Error(err)
			}
		}
		postUploadBridgeChannel := bridge.GetUploadTransBridgeInstance()
		if postUploadBridgeChannel != nil {
			tmpData, _ := postUploadBridgeChannel.PopAll()
			tmpDBData := make([]*dbData[model.BridgeDataModel], len(tmpData))
			for i, data := range tmpData {
				tmpDBData[i] = &dbData[model.BridgeDataModel]{Tag: BridgeDataTag, Data: data}
			}
			if err := InsertBatchDbData(tmpDBData); err != nil {
				logrus.Error(err)
			}
		}
	}
}

func receiveDataToBridge[T model.BridgeDataModel](data []*dbData[T]) {
	ia := bridge.GetPreUploadTransBridgeInstance()
	ib := bridge.GetUploadTransBridgeInstance()
	iaList := make([]*model.ScraperPreUploadFileDataModel, 0)
	ibList := make([]*model.ScraperPostUploadFileDataModel, 0)
	for _, d := range data {
		v := reflect.ValueOf(d.Data)
		switch v.Elem().Interface().(type) {
		case model.ScraperPreUploadFileDataModel:
			iaList = append(iaList, v.Interface().(*model.ScraperPreUploadFileDataModel))
		case model.ScraperPostUploadFileDataModel:
			ibList = append(ibList, v.Interface().(*model.ScraperPostUploadFileDataModel))
		case model.PreTransformDataModel:
			logrus.Warning("Wait to implement")
		case model.PostTransformDataModel:
			logrus.Warning("Wait to implement")
		default:
			logrus.Error(fmt.Sprintf("Unknown type: %s", v.Elem().Type().Name()))
		}
	}
	if len(iaList) > 0 {
		if err := ia.Insert(iaList); err != nil {
			logrus.Error(err)
		}
	}
	if len(ibList) > 0 {
		if err := ib.Insert(ibList); err != nil {
			logrus.Error(err)
		}
	}
}
