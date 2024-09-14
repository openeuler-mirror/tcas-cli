package initconfig

import (
	"encoding/json"
	"fmt"
	"gitee.com/nanhu-lab/tcas-cli/manager"
	"gitee.com/nanhu-lab/tcas-cli/utils/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path"
)

const (
	ConfigFile = ".tcasctl/config.json"
)

type TCASConfig struct {
	TCASInfos map[string]*manager.TCASInfo `json:"configs"`
}

// Cmd represents the ca command
var Cmd = &cobra.Command{
	Use:                        "init",
	Short:                      "init config",
	Long:                       "generate the config file",
	SuggestionsMinimumDistance: 1,
	DisableSuggestions:         false,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")

		homedir, err := os.UserHomeDir()
		if err != nil {
			logrus.Errorf("home dir not find, error: %s", err)
			return
		}
		configPath := path.Join(homedir, ConfigFile)
		configs := new(TCASConfig)
		if file.IsFile(configPath) == false {
			logrus.Debugf("config file does not exist, generate new config file")
			err := os.MkdirAll(path.Dir(configPath), os.ModeDir)
			if err != nil {
				logrus.Errorf("creat config dir failed, error: %s", err)
				return
			}
			configs.TCASInfos = map[string]*manager.TCASInfo{url: {
				APIEndpoint: url,
			}}
		} else {
			configs, err = ReadTCASConfigs(configPath)
			if err != nil {
				logrus.Errorf("read tcas configs failed, error: %s", err)
				return
			}
			if _, ok := configs.TCASInfos[url]; !ok {
				configs.TCASInfos[url] = &manager.TCASInfo{
					APIEndpoint: url,
				}
			} else {
				configs = nil
			}
		}

		if configs != nil {
			err = WriteTCASConfigToFile(configs, configPath)
			if err != nil {
				logrus.Errorf("write config faield, error: %s", err)
				return
			}
			fmt.Printf("init succuessful, config path: %s\n", configPath)
		} else {
			fmt.Println("config is already init")
		}
	},
}

func WriteTCASConfigToFile(configs *TCASConfig, savePath string) error {
	configsJson, err := json.MarshalIndent(configs, "", " ")
	if err != nil {
		return fmt.Errorf("marshal config to json failed, error: %s", err)
	}

	err = os.WriteFile(savePath, configsJson, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func ReadTCASConfigs(configPath string) (*TCASConfig, error) {
	configByte, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file failed, error: %s", err)
	}

	configs := new(TCASConfig)
	err = json.Unmarshal(configByte, configs)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config faield, error: %s", err)
	}

	return configs, nil
}

func ReadTCASInfoFromConfig(url string) (*manager.TCASInfo, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("home dir not find, error: %s", err)
	}
	tcasConfig, err := ReadTCASConfigs(path.Join(homedir, ConfigFile))
	if err != nil {
		return nil, err
	}
	if tcasInfo, ok := tcasConfig.TCASInfos[url]; ok {
		return tcasInfo, nil
	} else {
		return nil, fmt.Errorf("login info does not exist")
	}

}

func init() {
	Cmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
}
