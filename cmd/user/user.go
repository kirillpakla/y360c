// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package user

import (
	"github.com/foxsoft2005/y360c/cmd/user/alias"
	"github.com/foxsoft2005/y360c/cmd/user/contact"
	"github.com/foxsoft2005/y360c/cmd/user/mfa"
	"github.com/foxsoft2005/y360c/cmd/user/rule"
	"github.com/spf13/cobra"
)

var (
	orgId     int
	token     string
	userId    string
	userEmail string
	asRaw     bool
)

var Cmd = &cobra.Command{
	Use:   "user",
	Short: "Manage Y360 users",
	Long: `The command provides access to Y360 user management.
Cannot be executed directly, use one of the available sub-commands.
Y360_TOKEN and Y360_ORG_ID are read from the environment. Y360_API_BASE is read from .env.`,
	Example: `  y360c user add --help
  y360c user set --help
  y360c user ls --email person@example.com`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("user called")
	//	},
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(infoCmd)
	Cmd.AddCommand(mfa.Cmd)
	Cmd.AddCommand(alias.Cmd)
	Cmd.AddCommand(rule.Cmd)
	Cmd.AddCommand(contact.Cmd)
	Cmd.AddCommand(setCmd)
	Cmd.AddCommand(addCmd)
	Cmd.AddCommand(rmCmd)
	Cmd.AddCommand(senderinfoCmd)
}
