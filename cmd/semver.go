package main

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

var (
	GitCommit = "unknown"
	BuildDate = "unknown"
	Version   = "unreleased"
)

func printVersion() {
	fmt.Printf("Version:\t %s\n", Version)
	fmt.Printf("Git commit:\t %s\n", GitCommit)
	fmt.Printf("Date:\t\t %s\n", BuildDate)
}

func printUsage() {
	fmt.Printf(`Usage: semver <test-string> <semver-test-constraint>
Example:
  Test if a version is higher than 1.0.0:
  semver "1.0.1" ">1.0.0"

  Test if a version is less than 2.0.0 or higher than 3.0.0:
  semver "1.0.1" "<2.0.0 || >3.0.0"

  You can also use x, X or * as wildcard character like this:
  semver "1.0.0" "1.x"

  Or like this
  semver "2.1.9" ">2.1.x"
`)
}

func newRootCommand() *cobra.Command {
	var (
		showVersion bool
		quietMode   bool
	)

	var cmd = &cobra.Command{
		Use:          "semver-cli",
		Short:        "semver-cli is a helper for working around semver",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if showVersion {
				printVersion()
				return nil
			}

			if len(args) != 2 {
				if !quietMode {
					fmt.Fprintln(os.Stderr, "semver required 2 parameters.")
					printUsage()
				}
				os.Exit(128)
			}

			testVersion := args[0]
			inputConstraint := args[1]

			constraint, err := semver.NewConstraint(inputConstraint)
			if err != nil {
				if !quietMode {
					fmt.Fprintf(os.Stderr, "constraint %v is not a valid constraint.", inputConstraint)
					fmt.Fprintf(os.Stderr, "Error: %q", err)
				}
				os.Exit(128)
			}

			version, err := semver.NewVersion(testVersion)
			if err != nil {
				if !quietMode {
					fmt.Fprintf(os.Stderr, "%v is not a valid semver.", testVersion)
					fmt.Fprintf(os.Stderr, "Error: %q", err)
				}
				os.Exit(128)
			}

			valid, _ := constraint.Validate(version)

			if valid {
				if !quietMode {
					fmt.Printf("%v satisfies the semver constraint.\n", testVersion)
				}
				os.Exit(0)
			} else {
				if !quietMode {
					fmt.Printf("%v does not satisfies the constraint!\n", testVersion)
				}
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Show version info.")
	cmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "Quiet mode.")
	return cmd
}

func main() {
	cmd := newRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
