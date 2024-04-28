package sqlite

import (
	"bytes"
	"encoding/gob"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/model"
)

func init() {
	gob.Register(&model.ScraperPreUploadFileDataModel{})
	gob.Register(&model.ScraperPostUploadFileDataModel{})
	gob.Register(&model.PreTransformDataModel{})
	gob.Register(&model.PostTransformDataModel{})
}

// decodeData decodes the byte slice into a dbData
func decodeData(data []byte) (*dbData, error) {
	buffer := bytes.NewReader(data)
	dataDecoder := gob.NewDecoder(buffer)
	var result dbData
	if err := dataDecoder.Decode(&result); err != nil {
		return &dbData{}, log.ErrorWrap(err)
	}
	return &result, nil
}

// decodeDataBatch decodes a slice of byte slices into a slice of dbData
func decodeDataBatch(data [][]byte) ([]*dbData, error) {
	var results []*dbData
	for _, d := range data {
		buffer := bytes.NewReader(d)
		dataDecoder := gob.NewDecoder(buffer)
		result := dbData{}
		if err := dataDecoder.Decode(&result); err != nil {
			return nil, log.ErrorWrap(err)
		}
		results = append(results, &result)
	}
	return results, nil
}
