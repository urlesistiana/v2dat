package unpack

import (
	"github.com/spf13/cobra"
	"github.com/urlesistiana/v2dat/cmd"
)

type UnpackArgs struct {
	outDir           string
	file             string
	filters          []string
	with_type_prefix bool
}

var unpack = &cobra.Command{
	Use:   "unpack",
	Short: "unpack geosite and geoip to text files.",
}

func init() {
	unpack.AddCommand(newGeoSiteCmd(), newGeoIPCmd())
	cmd.RootCmd.AddCommand(unpack)
}

func AddCommand(cmds ...*cobra.Command) {
	unpack.AddCommand(cmds...)
}
