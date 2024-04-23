package storage

import (
	"NekoImageWorkflowKitex/common"
	"NekoImageWorkflowKitex/uploadClient/model"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var configPath string
var configFileName = "NekoImageWorkflowClientConfig"
var configFileNameWithExtension = "NekoImageWorkflowClientConfig.json"
var loadConfigOnce sync.Once
var inCacheConfig *model.ClientConfig

func loadConfig(info *model.ClientConfig) error {
	var config model.ConfigWrapper
	exe, err := os.Executable()
	configPath = filepath.Dir(exe)
	if err != nil {
		logrus.Error("Error getting current directory: %s\n", err)
		return err
	}
	if _, err := os.Stat(filepath.Join(configPath, configFileNameWithExtension)); os.IsNotExist(err) {
		CreateConfig()
	} else {
		viper.SetConfigName(configFileName)
		viper.AddConfigPath(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			logrus.Error("Error reading config file, ", err)
			return err
		}
		err = viper.Unmarshal(&config)
		if err != nil {
			logrus.Error("Error unmarshalling config file, ", err)
			return err
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
	viper.Set("ClientConfig", model.ClientConfig{
		ClientID:              "example-id",
		ClientName:            "example-name",
		DestServiceName:       "example-service",
		ClientRegisterAddress: "https://example.com/register",
		ConsulAddress:         "https://example-consul.com",
		PostUploadPeriod:      300,
		ScraperList:           []common.ScraperType{common.LocalScraperType, common.APIScraperType},
		ScraperConfig: model.ScraperConfig{
			LocalScraperConfig: model.LocalScraperConfig{
				WatchFolders: []string{"/path/to/watch/folder1", "/path/to/watch/folder2"},
			},
			APIScraperConfig: model.APIScraperConfig{
				APIScraperSource: []model.APIScraperSourceConfig{
					{
						APIAddress:           "https://example.com/api",
						ParserJavaScriptFile: "example-parser.js",
					},
				},
			},
		},
	})
	err := viper.SafeWriteConfig()
	if err != nil {
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

func GetConfig() *model.ClientConfig {
	loadConfigOnce.Do(func() {
		var config model.ClientConfig
		err := loadConfig(&config)
		if err != nil {
			inCacheConfig = &model.ClientConfig{}
		} else {
			inCacheConfig = &config
		}
	})
	return inCacheConfig
}
