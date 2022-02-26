package cmd

import (
	"jakubenglicky/kubessh/cmd/root"
	"github.com/spf13/cobra"
)
func Execute() {
	cobra.CheckErr(root.RootCmd.Execute())
}