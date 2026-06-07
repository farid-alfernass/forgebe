package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "forgebe",
		Short: "ForgeBE is an AI-agnostic backend engineering framework",
		Long:  "ForgeBE helps backend engineers keep AI-assisted delivery consistent, safe, local-first, and project-aware.",
	}

	// Persistent flags available to all commands
	cmd.PersistentFlags().BoolP("json", "j", false, "Output in JSON format")

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newScanCmd())
	cmd.AddCommand(newDoctorCmd())
	cmd.AddCommand(newProfileCmd())
	cmd.AddCommand(newExportCmd())
	cmd.AddCommand(newBriefCmd())
	cmd.AddCommand(newImportCmd())
	return cmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
