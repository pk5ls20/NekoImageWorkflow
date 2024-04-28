package sqlite

import (
	"bytes"
	"encoding/gob"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/model"
	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&model.ScraperPreUploadFileDataModel{})
	gob.Register(&model.ScraperPostUploadFileDataModel{})
	gob.Register(&model.PreTransformDataModel{})
	gob.Register(&model.PostTransformDataModel{})
}

// encodeData encodes the dbData into a byte slice
func encodeData(data *dbData) ([]byte, error) {
	buffer := new(bytes.Buffer)
	dataEncoder := gob.NewEncoder(buffer)
	if err := dataEncoder.Encode(data); err != nil {
		return nil, log.ErrorWrap(err)
	}
	result := buffer.Bytes()
	return result, nil
}

// encodeDataBatch encodes a slice of dbData into a slice of byte slices.
func encodeDataBatch(data *[]dbData) ([][]byte, error) {
	var results [][]byte
	for _, d := range *data {
		buffer := new(bytes.Buffer)
		dataEncoder := gob.NewEncoder(buffer)
		if err := dataEncoder.Encode(d); err != nil {
			logrus.Error("Failed to encode data: ", err)
			return nil, log.ErrorWrap(err)
		}
		result := buffer.Bytes()
		results = append(results, result)
	}
	return results, nil
}
