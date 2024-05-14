package config

import (
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var configPath string
var loadConfigOnce sync.Once
var configImpl *ClientConfig
var configFileName = "NekoImageWorkflowClientConfig"
var configFileNameWithExtension = "NekoImageWorkflowClientConfig.json"

func loadConfig(info *ClientConfig) error {
	var config ConfigWrapper
	exe, exeErr := os.Executable()
	if exeErr != nil {
		return log.ErrorWrap(exeErr)
	}
	configPath = filepath.Dir(exe)
	if _, fileErr := os.Stat(filepath.Join(configPath, configFileNameWithExtension)); os.IsNotExist(fileErr) {
		CreateConfig()
	} else {
		viper.SetConfigName(configFileName)
		viper.AddConfigPath(configPath)
		if err := viper.ReadInConfig(); err != nil {
			return log.ErrorWrap(err)
		}
		if err := viper.Unmarshal(&config); err != nil {
			return log.ErrorWrap(err)
		}
		*info = config.ClientConfig
	}
	return nil
}

func CreateConfig() {
	logrus.Warn("Config file not found, creating new one.")
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")
	// TODO: update the outdated config fit for the new code in main.go
	viper.Set("ClientConfig", ClientConfig{
		ClientRegisterAddress: "https://example.com/register",
		ConsulAddress:         "https://example-consul.com",
		PostUploadPeriod:      300,
		Scraper: map[commonModel.ScraperType][]Scraper{
			commonModel.LocalScraperType: {
				LocalScraperConfig{
					Enable:       true,
					WatchFolders: []string{"/path/to/watch/folder1", "/path/to/watch/folder2"},
				},
			},
			commonModel.APIScraperType: {
				APIScraperConfig{
					Enable: true,
					APIScraperSource: []APIScraperSourceConfig{
						{
							APIAddress:           "https://example-api.com",
							ParserJavaScriptFile: "example-parser.js",
							OptionalHeaders: map[string]string{
								"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
									"AppleWebKit/537.36 (KHTML, like Gecko) " +
									"Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0",
							},
							OptionalCookies: map[string]string{},
						},
					},
				},
			},
		},
	})
	if err := viper.SafeWriteConfig(); err != nil {
		var configFileAlreadyExistsError viper.ConfigFileAlreadyExistsError
		if errors.As(err, &configFileAlreadyExistsError) {
			logrus.Error("In CreateConfig(), Config file already exists.")
		}
	}
	logrus.Warn("Restart the program to load the new config.")
	// Hang the program to prevent it from exiting
	for {
		time.Sleep(114514 * time.Second)
	}
}

func GetConfig() *ClientConfig {
	loadConfigOnce.Do(func() {
		var config ClientConfig
		if err := loadConfig(&config); err != nil {
			logrus.Fatal("Failed to load config:", err)
		}
		configImpl = &config
	})
	return configImpl
}
