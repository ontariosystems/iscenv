package cmd

import (
	"github.com/ontariosystems/iscenv/v3/internal/cmd/flags"
	"github.com/spf13/cobra"
)

const userFlag = "user"

func addDockerUserFlags(cmd *cobra.Command) {
	flags.AddConfigFlagP(cmd, userFlag, "u", "", "User to run the command")
}
