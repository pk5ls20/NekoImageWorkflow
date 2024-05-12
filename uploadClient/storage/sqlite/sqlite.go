package sqlite

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	fileQueue "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/queue"
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
// InitSqlite Initialises the database and writes the residual buffered data within the database to fileQueue
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
		receiveData, err := FindDbDataModelByTag(QueueDataTag)
		if err != nil {
			logrus.Fatal(err)
		}
		tmpBytes := make([][]byte, len(receiveData))
		for i, data := range receiveData {
			tmpBytes[i] = data.Data
		}
		decodedReceiveData, err := decodeDataBatch[model.FileDataModel](tmpBytes)
		if err != nil {
			logrus.Fatal(err)
		}
		receiveDataToQueue[model.FileDataModel](decodedReceiveData)
		// 3. delete data from sqlite
		if err_ := DeleteDbDataByTag(QueueDataTag); err_ != nil {
			logrus.Fatal(err_)
		}
	})
}

// CloseSqlite TODO: impl tx
// CloseSqlite Securely close the database and write the contents of fileQueue
func CloseSqlite() {
	if dbInstance != nil {
		pushQueueData()
	}
}

func LoadClientUUID() (uuid.UUID, error) {
	receiveData, err := FindDbDataModelByTag(UUIDTag)
	if err != nil {
		logrus.Fatal(err)
	}
	if len(receiveData) == 0 {
		logrus.Warning("Client uuid not found, generating new one")
		_uuid := uuid.New()
		data := &dbData[uuid.UUID]{Tag: UUIDTag, Data: _uuid}
		if _err := InsertDbData[uuid.UUID](data); _err != nil {
			logrus.Fatal(err)
		}
		return _uuid, nil
	}
	if len(receiveData) > 0 {
		logrus.Debug("Client uuid found")
		_decodeData, _err := decodeData[uuid.UUID](receiveData[0].Data)
		if _err != nil {
			logrus.Fatal(err)
		}
		return _decodeData.Data, nil
	}
	return uuid.UUID{}, nil
}

// TODO: maybe we can move this out of sqlite.go, Consider dynamic registration functions
func receiveDataToQueue[T model.FileDataModel](data []*dbData[T]) {
	ia := fileQueue.GetPreUploadQueueInstance()
	ib := fileQueue.GetUploadQueueInstance()
	iaList := make([]*model.PreUploadFileDataModel, 0)
	ibList := make([]*model.UploadFileDataModel, 0)
	for _, d := range data {
		v := reflect.ValueOf(d.Data)
		switch v.Elem().Interface().(type) {
		case model.PreUploadFileDataModel:
			iaList = append(iaList, v.Interface().(*model.PreUploadFileDataModel))
		case model.UploadFileDataModel:
			ibList = append(ibList, v.Interface().(*model.UploadFileDataModel))
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

func pushQueueData() {
	// TODO: maybe we can move this out of sqlite.go, Consider dynamic registration functions
	// PreUploadFileDataModel
	preUploadQueue := fileQueue.GetPreUploadQueueInstance()
	if preUploadQueue != nil {
		tmpData, _ := preUploadQueue.PopAll()
		tmpDBData := make([]*dbData[model.FileDataModel], len(tmpData))
		for i, data := range tmpData {
			tmpDBData[i] = &dbData[model.FileDataModel]{Tag: QueueDataTag, Data: data}
		}
		if err := InsertBatchDbData(tmpDBData); err != nil {
			logrus.Error(err)
		}
	}
	// UploadFileDataModel
	postUploadQueue := fileQueue.GetUploadQueueInstance()
	if postUploadQueue != nil {
		tmpData, _ := postUploadQueue.PopAll()
		tmpDBData := make([]*dbData[model.FileDataModel], len(tmpData))
		for i, data := range tmpData {
			tmpDBData[i] = &dbData[model.FileDataModel]{Tag: QueueDataTag, Data: data}
		}
		if err := InsertBatchDbData(tmpDBData); err != nil {
			logrus.Error(err)
		}
	}
}
