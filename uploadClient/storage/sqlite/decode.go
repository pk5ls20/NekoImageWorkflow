package sqlite

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
)

func init() {
	// TODO: auto register, or add more!
	gob.Register(&uuid.UUID{})
	gob.Register(&clientModel.PreUploadFileDataModel{})
	gob.Register(&clientModel.UploadFileDataModel{})
}

// decodeData decodes the byte slice into a dbData
func decodeData[T dbDataModel](data []byte) (*dbData[T], error) {
	buffer := bytes.NewReader(data)
	dataDecoder := gob.NewDecoder(buffer)
	var result dbData[T]
	if err := dataDecoder.Decode(&result); err != nil {
		return &dbData[T]{}, log.ErrorWrap(err)
	}
	return &result, nil
}

// decodeDataBatch decodes a slice of byte slices into a slice of dbData
func decodeDataBatch[T dbDataModel](data [][]byte) ([]*dbData[T], error) {
	var results []*dbData[T]
	for _, d := range data {
		buffer := bytes.NewReader(d)
		dataDecoder := gob.NewDecoder(buffer)
		result := dbData[T]{}
		if err := dataDecoder.Decode(&result); err != nil {
			return nil, log.ErrorWrap(err)
		}
		results = append(results, &result)
	}
	return results, nil
}
