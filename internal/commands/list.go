package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list your Discord application commands",
	Long:    "list the Discord application commands associated with a Discord application (global or guild)",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cfg == nil {
			return nil
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("error validating config: %w", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString(discordBotTokenEnv)
		if token == "" {
			return fmt.Errorf("must pass in Discord Bot token via the --discord-bot-token (-t) flag or DISCORD_BOT_TOKEN environment variable")
		}

		appId := viper.GetString(discordAppIdEnv)
		if appId == "" {
			return fmt.Errorf("must pass in Discord Application ID via the --discord-app-id (-a) flag or DISCORD_APP_ID environment variable")
		}

		guildId := viper.GetString(discordGuildIdEnv)

		session, err := discordgo.New("Bot " + token)
		if err != nil {
			return err
		}

		commands, err := session.ApplicationCommands(appId, guildId)
		if err != nil {
			return err
		}

		for _, command := range commands {
			fmt.Printf("%s (%s): %s\n", command.Name, command.ID, command.Description)
		}
		return nil
	},
}
