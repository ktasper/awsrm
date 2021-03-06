package cmd

import (
	"log"

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
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "", false, "Enables debug mode (Also enables verbose mode)")
	rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "Enables quiet mode, you will not be prompted for any input")
	rootCmd.PersistentFlags().StringVarP(&awsRegion, "region", "r", "us-east-1", "The AWS region you wish to use")
	rootCmd.PersistentFlags().StringVarP(&awsProfile, "profile", "p", "default", "The AWS profile you wish to use")
	rootCmd.PersistentFlags().BoolVarP(&dryRunMode, "dry-run", "", false, "Enables dry-run mode (Will not make any changes)")
	rootCmd.PersistentFlags().BoolVarP(&skipVpcCheck, "skip-vpc-check", "", false, "Skip VPC Check; Skip the check to see if a matching vpc exists in the region")

	// Enable verbose mode if debug flag is set
	if debugMode {
		verboseMode = true
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Here is vars for all the flags I want to use
var (
	verboseMode  bool
	debugMode    bool
	quietMode    bool
	dryRunMode   bool
	skipVpcCheck bool
	awsRegion    string
	awsProfile   string
)

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
