package login

import (
	"fmt"
	"gitee.com/nanhu-lab/tcas-cli/cmd/initconfig"
	"gitee.com/nanhu-lab/tcas-cli/manager"
	"gitee.com/nanhu-lab/tcas-cli/tees"
	"gitee.com/nanhu-lab/tcas-cli/utils/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"path"
)

// Cmd represents the ca command
var Cmd = &cobra.Command{
	Use:                        "login",
	Short:                      "login tcas",
	Long:                       "login tcas",
	SuggestionsMinimumDistance: 1,
	DisableSuggestions:         false,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		skipssl, _ := cmd.Flags().GetBool("skipssl")

		if username == "" {
			fmt.Print("Username: ")
			fmt.Scanln(&username)
		}

		if password == "" {
			fmt.Print("Password: ")
			bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				logrus.Errorf("read password failed, error: %s", err)
				return
			}
			password = string(bytePassword)
		} else {
			fmt.Println("WARNING! Using --password via the CLI is insecure")
		}

		var tcasConfig *manager.TCASInfo
		tcasConfig, err := initconfig.ReadTCASInfoFromConfig(url)
		if err != nil {
			// try to login use default config
			tcasConfig = new(manager.TCASInfo)
			tcasConfig.APIEndpoint = url
		}
		tcasConfig.SkipVerify = skipssl
		m, err := manager.New(tcasConfig, tees.GetCollectors())
		if err != nil {
			logrus.Errorf("create attest manager failed, error: %s", err)
			return
		}

		config, err := m.Login(username, password)
		if err != nil {
			logrus.Errorf("%s", err)
			return
		}

		homedir, err := os.UserHomeDir()
		if err != nil {
			logrus.Errorf("home dir not find, error: %s", err)
			return
		}
		configPath := path.Join(homedir, initconfig.ConfigFile)
		configs := new(initconfig.TCASConfig)
		if file.IsFile(configPath) == false {
			logrus.Debugf("config file does not exist, generate new config file")
			err := os.MkdirAll(path.Dir(configPath), os.ModeDir)
			if err != nil {
				logrus.Errorf("creat config dir faile, error: %s", err)
				return
			}
			config.APIEndpoint = url
			configs.TCASInfos = map[string]*manager.TCASInfo{url: config}
		} else {
			configs, err = initconfig.ReadTCASConfigs(configPath)
			if err != nil {
				logrus.Errorf("read tcas confings failed, error: %s", err)
				return
			}
			configs.TCASInfos[url] = config
		}
		err = initconfig.WriteTCASConfigToFile(configs, configPath)
		if err != nil {
			logrus.Errorf("write config failed, error: %s", err)
			return
		}
		fmt.Println("login successful")
	},
}

func init() {
	Cmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	Cmd.Flags().StringP("username", "n", "", "optional, tcas's username")
	Cmd.Flags().StringP("password", "p", "", "optional, tcas's password")
	Cmd.Flags().BoolP("skipssl", "s", false, "optional, skip ssl verify")
}
