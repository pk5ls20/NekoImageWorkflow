package sqlite

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/msgQueue"
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
// InitSqlite Initialises the database and writes the residual buffered data within the database to storageQueue
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
		receiveData, err := FindDbDataModelByTag(ActivateQueueDataTag)
		if err != nil {
			logrus.Fatal(err)
		}
		tmpBytes := make([][]byte, len(receiveData))
		for i, data := range receiveData {
			tmpBytes[i] = data.Data
		}
		decodedReceiveData, err := decodeDataBatch[clientModel.FileDataModel](tmpBytes)
		if err != nil {
			logrus.Fatal(err)
		}
		receiveDataToQueue[clientModel.FileDataModel](decodedReceiveData)
		// 3. delete data from sqlite
		if err_ := DeleteDbDataByTag(ActivateQueueDataTag); err_ != nil {
			logrus.Fatal(err_)
		}
	})
}

// CloseSqlite TODO: impl tx
// CloseSqlite Securely close the database and write the contents of storageQueue
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
func receiveDataToQueue[T clientModel.FileDataModel](data []*dbData[T]) {
	queue := msgQueue.NewMessageQueue()
	msgDataList := make([]*msgQueue.MsgQueueData, 0)
	for _, d := range data {
		v := reflect.ValueOf(d.Data)
		switch v.Elem().Interface().(type) {
		case clientModel.PreUploadFileDataModel:
			ori := v.Interface().(*clientModel.PreUploadFileDataModel)
			model := &msgQueue.MsgQueueData{
				MsgMetaData: msgQueue.MsgMetaData{
					UploadType: commonModel.PreUploadType,
					MsgMetaID: msgQueue.MsgMetaID{
						ScraperType: ori.ScraperType,
						ScraperID:   ori.GetScraperID(),
						MsgGroupID:  ori.MsgGroupID,
					},
				},
				FileMetaData: &clientModel.AnyFileMetaDataModel{
					PreUploadFileMetaDataModel: &ori.PreUploadFileMetaDataModel,
				},
			}
			msgDataList = append(msgDataList, model)
		case clientModel.UploadFileDataModel:
			ori := v.Interface().(*clientModel.UploadFileDataModel)
			model := &msgQueue.MsgQueueData{
				MsgMetaData: msgQueue.MsgMetaData{
					UploadType: commonModel.PostUploadType,
					MsgMetaID: msgQueue.MsgMetaID{
						ScraperType: ori.ScraperType,
						ScraperID:   ori.GetScraperID(),
						MsgGroupID:  ori.MsgGroupID,
					},
				},
				FileMetaData: &clientModel.AnyFileMetaDataModel{
					UploadFileMetaDataModel: &ori.UploadFileMetaDataModel,
				},
			}
			msgDataList = append(msgDataList, model)
		default:
			logrus.Error(fmt.Sprintf("Unknown type: %s", v.Elem().Type().Name()))
		}
	}
	if err := queue.AddElements(msgDataList); err != nil {
		logrus.Error(err)
	}
}

func transMsgQueueData(tag keyTag, queueData []*msgQueue.MsgQueueData) ([]*dbData[clientModel.FileDataModel], error) {
	tmpDBData := make([]*dbData[clientModel.FileDataModel], len(queueData))
	for idx, itm := range queueData {
		if itm.FileMetaData.PreUploadFileMetaDataModel != nil {
			tmpDBData[idx] = &dbData[clientModel.FileDataModel]{
				Tag: tag,
				Data: &clientModel.PreUploadFileDataModel{
					PreUploadFileMetaDataModel: *itm.FileMetaData.PreUploadFileMetaDataModel,
				},
			}
		} else if itm.FileMetaData.UploadFileMetaDataModel != nil {
			tmpDBData[idx] = &dbData[clientModel.FileDataModel]{
				Tag: tag,
				Data: &clientModel.UploadFileDataModel{
					UploadFileMetaDataModel: *itm.FileMetaData.UploadFileMetaDataModel,
				},
			}
		} else {
			return tmpDBData, errors.New("unknown type In pushQueueData")
		}
	}
	return tmpDBData, nil
}

func pushQueueData() {
	// TODO: maybe we can move this out of sqlite.go, Consider dynamic registration functions
	queue := msgQueue.NewMessageQueue()
	// activate msg queue
	activateMsgList, err := queue.PopAll(msgQueue.ActivateQueue)
	if err != nil {
		logrus.Error(err)
	}
	activateQueueData, err := transMsgQueueData(ActivateQueueDataTag, activateMsgList)
	if err != nil {
		logrus.Error(err)
	}
	if err = InsertBatchDbData[clientModel.FileDataModel](activateQueueData); err != nil {
		logrus.Error(err)
	}
	// dead msg queue
	deadMsgList, err := queue.PopAll(msgQueue.DeadQueue)
	if err != nil {
		logrus.Error(err)
	}
	deadQueueData, err := transMsgQueueData(DeadQueueDataTag, deadMsgList)
	if err != nil {
		logrus.Error(err)
	}
	if err = InsertBatchDbData[clientModel.FileDataModel](deadQueueData); err != nil {
		logrus.Error(err)
	}
}
