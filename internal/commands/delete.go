package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [command]",
	Aliases: []string{"rm", "del"},
	Short:   "delete Discord application command(s)",
	Long:    "delete Discord application command(s) associated with a Discord application (global or guild)",
	Args:    cobra.MinimumNArgs(1),
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

		for _, commandToDelete := range args {
			for _, command := range commands {
				if command.Name == commandToDelete {
					if guildId == "" {
						fmt.Printf("Deleting global command %s (%s)\n", command.Name, command.ID)
					} else {
						fmt.Printf("Deleting guild command %s (%s) from guild %s\n", command.Name, command.ID, command.GuildID)
					}

					err := session.ApplicationCommandDelete(appId, guildId, command.ID)
					if err != nil {
						if guildId == "" {
							fmt.Printf("Encountered error trying to delete global command %s (%s): %v\n", command.Name, command.ID, err)
						} else {
							fmt.Printf("Encountered error trying to delete guild command %s (%s) from guild %s: %v\n", command.Name, command.ID, command.GuildID, err)
						}
					} else {
						if guildId == "" {
							fmt.Printf("Sucessfully deleted global command %s (%s)\n", command.Name, command.ID)
						} else {
							fmt.Printf("Sucessfully deleted guild command %s (%s) from guild %s\n", command.Name, command.ID, command.GuildID)
						}
					}
					break
				}
			}

			if guildId == "" {
				fmt.Printf("Command '%s' not found within global commands\n", commandToDelete)
			} else {
				fmt.Printf("Command '%s' not found within guild commands for guild %s\n", commandToDelete, guildId)
			}
		}

		return nil
	},
}
