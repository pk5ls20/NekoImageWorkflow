package lifecycle

import (
	scraperModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	"github.com/sirupsen/logrus"
)

func StopScraper(scp []scraperModel.Scraper) {
	for _, scraperInstance := range scp {
		if err := scraperInstance.OnStop(); err != nil {
			logrus.Error("OnStop error:", err)
		}
	}
}
