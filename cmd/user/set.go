// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package user

import (
	"log"
	"os"
	"time"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	about        string
	birthday     string
	displayName  string
	externalId   string
	position     string
	isAdmin      helper.EnumYesNo
	nickname     string
	firstName    string
	lastName     string
	middleName   string
	gender       string
	language     string
	timezone     string
	isEnabled    helper.EnumYesNo
	departmentID int
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set an attributes of the user",
	Long: `Use this command to set an attributes (name, birthday, department, etc.) of the selected user.
"directory:write_users" permission is required (see Y360 help topics).`,
	Example: `  y360c user set --email person@example.com --department-id 42
  y360c user set --id user-100 --is-enabled no`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("user set is called")

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

		if birthday != "" {
			_, err := time.Parse(time.DateOnly, birthday)
			if err != nil {
				log.Fatalln("Birthdate cannot be parsed:", err)
			}
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

		patch := helper.PatchInput{
			About:       about,
			Birthday:    birthday,
			DisplayName: displayName,
			ExternalID:  externalId,
			Position:    position,
			Nickname:    nickname,
			Gender:      gender,
			Language:    language,
			Timezone:    timezone,
			First:       firstName,
			Last:        lastName,
			Middle:      middleName,
			IsAdmin:     isAdmin,
			IsEnabled:   isEnabled,
		}
		if cmd.Flags().Changed("department-id") {
			id := departmentID
			patch.DepartmentID = &id
		}
		data, err := helper.NewDirectoryClient(orgId, token).PatchEmployee(userId, patch)
		if err != nil {
			log.Fatalln("Unable to update user:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id})
		t.AppendRow(table.Row{"About", data.About})
		t.AppendRow(table.Row{"Birthday", data.Birthday})
		t.AppendRow(table.Row{"Display Name", data.DisplayName})
		t.AppendRow(table.Row{"External Id", data.ExternalId})
		t.AppendRow(table.Row{"Position", data.Position})
		t.AppendRow(table.Row{"Department Id", data.DepartmentId})
		t.AppendRow(table.Row{"Is Admin", data.IsAdmin})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	setCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	setCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	setCmd.Flags().StringVar(&userId, "id", "", "user id")
	setCmd.Flags().StringVar(&userEmail, "email", "", "user email address")

	setCmd.Flags().StringVar(&about, "about", "", "about an user")
	setCmd.Flags().StringVar(&birthday, "birthday", "", "user birthday (YYYY-MM-DD)")
	setCmd.Flags().StringVar(&displayName, "display-name", "", "user display name")
	setCmd.Flags().StringVar(&externalId, "external-id", "", "user external id")
	setCmd.Flags().StringVar(&position, "position", "", "user position")
	setCmd.Flags().Var(&isAdmin, "is-admin", "user has admin permissions (yes or no)")
	setCmd.Flags().StringVar(&nickname, "nickname", "", "employee login")
	setCmd.Flags().StringVar(&firstName, "first-name", "", "employee first name")
	setCmd.Flags().StringVar(&lastName, "last-name", "", "employee last name")
	setCmd.Flags().StringVar(&middleName, "middle-name", "", "employee middle name")
	setCmd.Flags().StringVar(&gender, "gender", "", "user gender")
	setCmd.Flags().StringVar(&language, "language", "", "user language")
	setCmd.Flags().StringVar(&timezone, "timezone", "", "user timezone")
	setCmd.Flags().Var(&isEnabled, "is-enabled", "user can sign in (yes or no)")
	setCmd.Flags().IntVar(&departmentID, "department-id", 0, "department id")

	setCmd.MarkFlagsOneRequired("id", "email")
	setCmd.MarkFlagsMutuallyExclusive("id", "email")

	setCmd.MarkFlagsOneRequired("about", "birthday", "display-name", "external-id", "position", "is-admin", "nickname", "first-name", "last-name", "middle-name", "gender", "language", "timezone", "is-enabled", "department-id")
}
