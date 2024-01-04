package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "bmctl",
		Short: "CLI to manage the bondMachine firmware images: build and load",
		Long:  `bondMachine cloud tools allow ....`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")

	buildCmd.PersistentFlags().StringP("bmfile", "f", "BMFile", "path to BMFile, relative to context")
	buildCmd.PersistentFlags().StringP("target", "t", "", "target image reference")
	buildCmd.PersistentFlags().StringP("platform", "p", "lattice/ice40/yosys", "platform name [vendor/board/variant]")

	buildCmd.MarkPersistentFlagRequired("target")

	loadCmd.PersistentFlags().StringP("cmd", "c", "", "custom command for loading firmware")
	loadCmd.PersistentFlags().StringP("device", "d", "", "device ID for target board")

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(loadCmd)
}

func initConfig() {
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := os.UserHomeDir()
	// 	cobra.CheckErr(err)

	// 	// Search config in home directory with name ".cobra" (without extension).
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigType("yaml")
	// 	viper.SetConfigName(".cobra")
	// }

	// viper.AutomaticEnv()

	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
}
