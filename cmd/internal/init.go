// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package internal

import (
	"errors"
	"log"
	"os"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the application",
	Long: `Use this command to check local settings.
Ordinary variables are read from .env. The access token is read from the environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Settings file is .env; the token is read from the environment")
		settings, err := helper.LoadSettings(".env", os.LookupEnv)
		if err != nil && !errors.Is(err, helper.ErrTokenMissing) {
			log.Fatalln("Incorrect settings:", err)
		}
		if settings.Token == "" {
			log.Print("Token: unset")
		} else {
			log.Print("Token: set")
		}
		log.Printf("Organization id: %s", helper.OrgIDStatus(settings.OrgID))
	},
}

func init() {

}
