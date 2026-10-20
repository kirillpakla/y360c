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
	Short: "Change employee fields",
	Long: `Use this command to change an existing employee.
"directory:write_users" permission is required (see Y360 help topics).

Select the employee with --email or --id. Pass only the fields to change.
A flag you omit is not sent, and that field stays as it is.

  --about           profile description
  --birthday        date as YYYY-MM-DD
  --display-name    name shown in the profile
  --external-id     external identifier
  --position        job title
  --department-id   department id; moves the employee
  --is-admin        administrator rights: yes or no
  --first-name      given name
  --last-name       family name
  --middle-name     patronymic
  --nickname        login
  --gender          gender value accepted by the directory
  --language        language code, for example ru
  --timezone        time zone, for example Europe/Moscow
  --is-enabled      sign-in allowed: yes or no

The result table shows id, description, birthday, display name, external id, job title, department and administrator state.
Name, login, gender, language, time zone and sign-in state change with their flags and are not printed in that table.`,
	Example: `  y360c user set --email person@example.com --position Engineer --department-id 42 --display-name "Person Example" --about "Directory account" --birthday 1990-01-31 --external-id ext-100
  y360c user set --email person@example.com --first-name Person --last-name Example --middle-name Middle --nickname person --gender unspecified --language ru --timezone Europe/Moscow
  y360c user set --id user-100 --is-enabled no --is-admin no`,
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
	setCmd.Flags().StringVar(&userId, "id", "", "user id of the employee to change")
	setCmd.Flags().StringVar(&userEmail, "email", "", "mailbox of the employee to change")

	setCmd.Flags().StringVar(&about, "about", "", "profile description")
	setCmd.Flags().StringVar(&birthday, "birthday", "", "birthday, YYYY-MM-DD")
	setCmd.Flags().StringVar(&displayName, "display-name", "", "name shown in the profile")
	setCmd.Flags().StringVar(&externalId, "external-id", "", "external identifier")
	setCmd.Flags().StringVar(&position, "position", "", "job title")
	setCmd.Flags().Var(&isAdmin, "is-admin", "administrator rights: yes or no")
	setCmd.Flags().StringVar(&nickname, "nickname", "", "login")
	setCmd.Flags().StringVar(&firstName, "first-name", "", "given name")
	setCmd.Flags().StringVar(&lastName, "last-name", "", "family name")
	setCmd.Flags().StringVar(&middleName, "middle-name", "", "patronymic")
	setCmd.Flags().StringVar(&gender, "gender", "", "gender value accepted by the directory")
	setCmd.Flags().StringVar(&language, "language", "", "language code, for example ru")
	setCmd.Flags().StringVar(&timezone, "timezone", "", "time zone, for example Europe/Moscow")
	setCmd.Flags().Var(&isEnabled, "is-enabled", "sign-in allowed: yes or no")
	setCmd.Flags().IntVar(&departmentID, "department-id", 0, "department id")

	setCmd.MarkFlagsOneRequired("id", "email")
	setCmd.MarkFlagsMutuallyExclusive("id", "email")

	setCmd.MarkFlagsOneRequired("about", "birthday", "display-name", "external-id", "position", "is-admin", "nickname", "first-name", "last-name", "middle-name", "gender", "language", "timezone", "is-enabled", "department-id")
}
