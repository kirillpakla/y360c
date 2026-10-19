// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package internal

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the application",
	Long: `Use this command to create / check the app configuration.
The config file will be automatically created if not found.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Config file is %s", viper.ConfigFileUsed())
		if viper.GetString("token") == "" {
			log.Printf("Token: unset")
		} else {
			log.Printf("Token: set")
		}
		log.Printf("Organization id: %d", viper.GetInt("org-id"))
	},
}

func init() {

}
