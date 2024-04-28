package lifecycle

import (
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper"
	"github.com/sirupsen/logrus"
)

func StartScraper(scp []scraper.ScraperInstance) {
	for _, scraperInstance := range scp {
		instance := scraperInstance
		go func() {
			logrus.Debug("go scraper: ", instance.GetType())
			if err := instance.PrepareData(); err != nil {
				logrus.Error("PrepareData error:", err)
			}
		}()
		go func() {
			logrus.Debug("go scraper: ", instance.GetType())
			if err := instance.ProcessData(); err != nil {
				logrus.Error("ProcessData error:", err)
			}
		}()
	}
}
