package commands

import (
	"github.com/BigPapaChas/dapp/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	cfgFile string
	cfg     *config.Config
	debug   bool
)

const (
	discordBotTokenEnv = "DISCORD_BOT_TOKEN"
	discordGuildIdEnv  = "DISCORD_GUILD_ID"
	discordAppIdEnv    = "DISCORD_APP_ID"
)

var rootCmd = &cobra.Command{
	Use:     "dapp",
	Short:   "dapp helps manage your Discord application commands",
	Version: "v0.0.1",
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gogok8s.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug messages")

	rootCmd.PersistentFlags().StringP("discord-bot-token", "t", "", "Discord Bot token for making requests (or setting DISCORD_BOT_TOKEN environment variable)")
	rootCmd.PersistentFlags().StringP("discord-guild-id", "g", "", "Discord Guild ID (or setting DISCORD_GUILD_ID environment variable)")
	rootCmd.PersistentFlags().StringP("discord-app-id", "a", "", "Discord Application ID (or setting DISCORD_APP_ID environment variable)")

	viper.BindPFlag(discordBotTokenEnv, rootCmd.PersistentFlags().Lookup("discord-bot-token"))
	viper.BindPFlag(discordGuildIdEnv, rootCmd.PersistentFlags().Lookup("discord-guild-id"))
	viper.BindPFlag(discordAppIdEnv, rootCmd.PersistentFlags().Lookup("discord-app-id"))

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dapp")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			cfg = config.NewConfig()
			cobra.CheckErr(viper.Unmarshal(cfg))
		} else {
			// Config file was found but another error was produced
			log.Println(err)
		}
	}
}

func Execute() error {
	return rootCmd.Execute()
}
