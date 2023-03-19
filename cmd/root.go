package cmd

import (
	"encoding/json"
	"os"

	"github.com/pete0emerson/lightbulb/lightbulb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var includeTags []string
var skipTags []string
var interactiveMode bool

func printJSON(obj interface{}) string {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lightbulb URL|FILE",
	Short: "A tool to execute annotated Markdown files",
	Long: `Lightbulb is a tool to execute annotated Markdown files. Leveraging HTML
comments in Markdown files, Lightbulb can create files and execute code blocks.
Lightbulb facilitates the testing of tutorial-style Markdown documentation.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("includeTags: %s", includeTags)
		log.Debugf("skipTags: %s", skipTags)
		content, err := lightbulb.LoadFromFile(args[0])
		if err != nil {
			log.Fatal(err)
		}
		blocks, err := lightbulb.Parse(string(content))
		if err != nil {
			log.Fatal(err)
		}

		err = lightbulb.Run(blocks)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lightbulb.yaml)")
	rootCmd.Flags().StringSliceVarP(&includeTags, "include-tag", "t", []string{"all"}, "Array stating the tags to be included")
	viper.BindPFlag("include-tag", rootCmd.Flags().Lookup("include-tag"))
	rootCmd.Flags().StringSliceVarP(&skipTags, "skip-tag", "x", []string{}, "Array stating the tags to be included")
	viper.BindPFlag("skip-tag", rootCmd.Flags().Lookup("skip-tag"))
	rootCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "Interactive mode")
	viper.BindPFlag("interactive", rootCmd.Flags().Lookup("interactive"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".lightbulb" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".lightbulb")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}
}
