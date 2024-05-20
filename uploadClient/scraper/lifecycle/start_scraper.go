package lifecycle

import (
	scraperModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	"github.com/sirupsen/logrus"
	"sync"
)

func StartScraper(scp []scraperModel.Scraper) {
	for _, scraperInstance := range scp {
		var wg sync.WaitGroup
		instance := scraperInstance
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := instance.OnStart(); err != nil {
				logrus.Error("OnStart error:", err)
			}
		}()
		wg.Wait()
		go func() {
			if err := instance.PrepareData(); err != nil {
				logrus.Error("PrepareData error:", err)
			}
		}()
		go func() {
			if err := instance.ProcessData(); err != nil {
				logrus.Error("ProcessData error:", err)
			}
		}()
	}
}
