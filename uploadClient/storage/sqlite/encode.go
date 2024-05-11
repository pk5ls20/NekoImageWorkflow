package sqlite

import (
	"bytes"
	"encoding/gob"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
)

func init() {
	// TODO: auto register, or add more!
	gob.Register(&clientModel.PreUploadFileDataModel{})
	gob.Register(&clientModel.UploadFileDataModel{})
}

// encodeData encodes the dbData into a byte slice
func encodeData[T dbDataModel](data *dbData[T]) ([]byte, error) {
	buffer := new(bytes.Buffer)
	dataEncoder := gob.NewEncoder(buffer)
	if err := dataEncoder.Encode(data); err != nil {
		return nil, log.ErrorWrap(err)
	}
	result := buffer.Bytes()
	return result, nil
}

// encodeDataBatch encodes a slice of dbData into a slice of byte slices.
func encodeDataBatch[T dbDataModel](data *[]dbData[T]) ([][]byte, error) {
	var results [][]byte
	for _, d := range *data {
		buffer := new(bytes.Buffer)
		dataEncoder := gob.NewEncoder(buffer)
		if err := dataEncoder.Encode(d); err != nil {
			return nil, log.ErrorWrap(err)
		}
		result := buffer.Bytes()
		results = append(results, result)
	}
	return results, nil
}
