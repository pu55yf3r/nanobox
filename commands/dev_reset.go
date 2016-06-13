package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// DevResetCmd ...
	DevResetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Resets the dev VM registry.",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			// TODO: Take an extra arguement and decide what we want to reset
			handleError(processor.Run("dev_reset", processor.DefaultConfig))
		},
	}
)