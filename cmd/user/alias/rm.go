// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package alias

import (
	"log"
	"os"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove an alias of the user mailbox",
	Long: `Use this command to remove an alias for the user mailbox.
"directory:write_users" permission is required (see Y360 help topics).`,
	Example: `  y360c user alias rm --email person@example.com --alias alias@example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user alias rm called")

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

		if !helper.Confirm("Do you REALLY want to DELETE the selected entity (y[es]|no)?") {
			log.Fatalln("Aborted by the user")
		}

		if userEmail != "" {
			us, err := helper.GetUserByEmail(orgId, token, userEmail)
			if err != nil {
				log.Fatalln("Failed to get user by email", err)
			}

			if us == nil {
				log.Fatalf("User (mailbox) %s does not found", userEmail)
			}

			userId = us.Id
		}

		data, err := helper.NewDirectoryClient(orgId, token).RemoveAlias(userId, alias)
		if err != nil {
			log.Fatalln("Unable to remove alias:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Alias", "Removed"})
		t.AppendRow(table.Row{data.Alias, data.Removed})
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	rmCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	rmCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	rmCmd.Flags().StringVar(&userId, "id", "", "user id")
	rmCmd.Flags().StringVar(&userEmail, "email", "", "user email address")
	rmCmd.Flags().StringVar(&alias, "alias", "", "alias to be deleted")

	rmCmd.MarkFlagsOneRequired("id", "email")
	rmCmd.MarkFlagsMutuallyExclusive("id", "email")

	err := rmCmd.MarkFlagRequired("alias")
	if err != nil {
		log.Fatalln("Error marking flag as required:", err)
	}
}
