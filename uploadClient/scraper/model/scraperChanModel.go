package model

import "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"

// ScraperChanMap stores the UploadFileDataModel that can be further processed
// scraper-PreUpload --(*model.PreUploadFileDataModel, all)--> client --(*model.PreUploadFileDataModel, need upload)-->
// scraper-Upload --(*model.UploadFileDataModel, all)--> client
// TODO: after impl messageQueue, maybe we don't need this anymore
type ScraperChanMap map[string]chan *model.PreUploadFileDataModel
