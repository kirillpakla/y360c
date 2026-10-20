// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package alias

import (
	"log"
	"os"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Get a list of the user mailbox aliases",
	Long: `Use this command to get a list of the user mailbox aliases.
"directory:read_users" permission is required (see Y360 help topics).`,
	Example: `  y360c user alias ls --email person@example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user alias ls called")

		if token == "" {
			t, err := helper.GetToken()
			if err != nil {
				log.Fatalln("Incorrect settings:", err)
			}
			token = t
		}

		if orgId == 0 {
			t, err := helper.GetOrgId()
			if err != nil {
				log.Fatalln("Incorrect settings:", err)
			}
			orgId = t
		}

		var data *model.User
		if userEmail != "" {
			found, err := helper.GetUserByEmail(orgId, token, userEmail)
			if err != nil {
				log.Fatalln("Failed to get user by email", err)
			}
			if found == nil {
				log.Fatalf("User (mailbox) %s does not found", userEmail)
			}
			data = found
		} else {
			found, err := helper.NewDirectoryClient(orgId, token).ListAliases(userId)
			if err != nil {
				log.Fatalln("Unable to list aliases:", err)
			}
			data = found
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"email", "user id", "alias"})
		if len(data.Aliases) == 0 {
			t.AppendRow(table.Row{data.Email, data.Id, ""})
		}
		for _, alias := range data.Aliases {
			t.AppendRow(table.Row{data.Email, data.Id, alias})
		}
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	lsCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	lsCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	lsCmd.Flags().StringVar(&userId, "id", "", "user id")
	lsCmd.Flags().StringVar(&userEmail, "email", "", "user email address")

	lsCmd.MarkFlagsOneRequired("id", "email")
	lsCmd.MarkFlagsMutuallyExclusive("id", "email")
}
