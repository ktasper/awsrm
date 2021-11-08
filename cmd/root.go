package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awsrm",
	Short: "Easily remove aws resources",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// If the user just types "awsrm" show the usage screen
		err := cmd.Usage()
		if err != nil {
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(&verboseMode, "verbose", "v", false, "Enables verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "Enables quiet mode, you will not be prompted for any input")
	rootCmd.PersistentFlags().StringVarP(&awsRegion, "region", "r", "us-east-1", "The AWS region you wish to use")
	rootCmd.PersistentFlags().StringVarP(&awsProfile, "profile", "p", "default", "The AWS profile you wish to use")
	rootCmd.PersistentFlags().BoolVarP(&dryRunMode, "dry-run", "", false, "Enables dry-run mode (Will not make any changes)")
	rootCmd.PersistentFlags().BoolVarP(&vpcSafeMode, "safe", "", true, "VPC Safe mode; Will delete current action ONLY if a matching vpc in the same region cannot be found")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Here is vars for all the flags I want to use
var (
	verboseMode bool
	quietMode   bool
	dryRunMode  bool
	vpcSafeMode bool
	awsRegion   string
	awsProfile  string
)

// Error handling
func exitErrorf(msg string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, msg+"\n", args...)
	if err != nil {
		return
	}
	os.Exit(1)
}

// Fancy prompt for users to get a bool value
func yesNo() bool {
	prompt := promptui.Select{
		Label: "Select[Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}
