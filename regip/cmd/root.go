package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/gookit/color"
	"regip"
	"regip/client"
)

var cfgFile string

// Used for CLI options
const DEFAULT_CLIENT_ADDR = "ws://localhost:8089/ws"

var clientAddr string
var logQuiet bool
var onlyID bool

func CreateLogger(name string, clr color.Color) *regip.Logger {
	if logQuiet || onlyID {
		return regip.NewLogger(ioutil.Discard)
	} else {
		return regip.NewLogger(os.Stdout).Tag(name, clr)

	}
}
func CreateClient(l *regip.Logger) (*client.Client, bool) {
	c, err := client.NewClient(clientAddr, l)
	if err != nil {
		l.Error("couldn't init client -- ", err)
		return nil, false
	}
	return c, true
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "regip",
	Short: "regip is the all-in-one binary for you name indexing needs! XD",
	Long: `regip is the program behind the Regex Indexing Project.
It includes the server, commands to import names, a CLI, and more!`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&clientAddr, "addr", "a", DEFAULT_CLIENT_ADDR, "address of regip server")
	rootCmd.PersistentFlags().BoolVarP(&logQuiet, "quiet", "q", false, "quiet logging")
	rootCmd.PersistentFlags().BoolVarP(&onlyID, "ids", "i", false, "only output")
	//rootCmd.MarkPersistentFlagRequired("db")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".regip" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".regip")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
