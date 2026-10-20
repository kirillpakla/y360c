// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package alias

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an alias for the user mailbox",
	Long: `Use this command to add an alias for the user mailbox.
"directory:write_users" permission is required (see Y360 help topics).`,
	Example: `  y360c user alias add --email person@example.com --alias alias@example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user alias add called")

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

		data, err := helper.NewDirectoryClient(orgId, token).AddAlias(userId, alias)
		if err != nil {
			log.Fatalln("Unable to add alias:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.ID})
		t.AppendRow(table.Row{"Name", strings.TrimSpace(fmt.Sprintf("%s %s %s", data.Name.Last, data.Name.First, data.Name.Middle))})
		t.AppendRow(table.Row{"Aliases", strings.Join(data.Aliases, ",")})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()

	},
}

func init() {
	addCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	addCmd.Flags().StringVar(&userId, "id", "", "user id")
	addCmd.Flags().StringVar(&userEmail, "email", "", "user email address")
	addCmd.Flags().StringVar(&alias, "alias", "", "mailbox alias")

	addCmd.MarkFlagsOneRequired("id", "email")
	addCmd.MarkFlagsMutuallyExclusive("id", "email")

	err := addCmd.MarkFlagRequired("alias")
	if err != nil {
		log.Fatalln("Error marking flag as required:", err)
	}
}
