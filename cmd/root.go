/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/foxsoft2005/y360c/cmd/auth"
	"github.com/foxsoft2005/y360c/cmd/dept"
	"github.com/foxsoft2005/y360c/cmd/disk"
	"github.com/foxsoft2005/y360c/cmd/dns"
	"github.com/foxsoft2005/y360c/cmd/domain"
	"github.com/foxsoft2005/y360c/cmd/group"
	"github.com/foxsoft2005/y360c/cmd/internal"
	"github.com/foxsoft2005/y360c/cmd/mailbox"
	"github.com/foxsoft2005/y360c/cmd/mfa"
	"github.com/foxsoft2005/y360c/cmd/org"
	"github.com/foxsoft2005/y360c/cmd/user"
	"github.com/foxsoft2005/y360c/cmd/whitelist"
	"github.com/foxsoft2005/y360c/helper"
)

var rootCmd = &cobra.Command{
	Use:   "y360c",
	Short: "Yandex360 cli",
	Long: `This app allows you to use the Yandex360 API via cli with some useful features.
Y360_TOKEN and Y360_ORG_ID are read from the environment. Y360_API_BASE is read from .env.`,
	Example: `  y360c user --help
  y360c user add --help
  y360c dept ls --help`,
	Version: "1.0.0-beta.47",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(auth.Cmd)
	rootCmd.AddCommand(mfa.Cmd)
	rootCmd.AddCommand(org.Cmd)
	rootCmd.AddCommand(dept.DeptCmd)
	rootCmd.AddCommand(user.Cmd)
	rootCmd.AddCommand(disk.Cmd)
	rootCmd.AddCommand(whitelist.Cmd)
	rootCmd.AddCommand(internal.InitCmd)
	rootCmd.AddCommand(group.GroupCmd)
	rootCmd.AddCommand(dns.Cmd)
	rootCmd.AddCommand(domain.Cmd)
	rootCmd.AddCommand(internal.Cmd)
	rootCmd.AddCommand(mailbox.Cmd)
}

func initConfig() {
	settings, err := helper.LoadSettings(".env", os.LookupEnv)
	if err != nil && !errors.Is(err, helper.ErrTokenMissing) {
		log.Fatalln("Incorrect settings:", err)
	}
	helper.ApplySettings(settings)
}
