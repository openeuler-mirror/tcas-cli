package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version string
var Cmd = &cobra.Command{
	Use:                "version",
	Args:               cobra.ArbitraryArgs,
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	Short:              "the version of tcasctl",
	Long:               "the version of tcasctl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version is %s\n", version)
	},
}
