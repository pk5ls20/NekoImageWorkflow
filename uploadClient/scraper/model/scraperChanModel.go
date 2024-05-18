package model

import "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"

// ScraperChanMap stores the UploadFileDataModel that can be further processed
// scraper-PreUpload --(*model.PreUploadFileDataModel, all)--> client --(*model.PreUploadFileDataModel, need upload)-->
// scraper-Upload --(*model.UploadFileDataModel, all)--> client
type ScraperChanMap map[int]chan *model.PreUploadFileDataModel
