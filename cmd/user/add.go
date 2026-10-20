package user

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	createNickname     string
	createDepartmentID int
	createFirst        string
	createLast         string
	createMiddle       string
	createAbout        string
	createBirthday     string
	createDisplayName  string
	createExternalID   string
	createGender       string
	createPosition     string
	createLanguage     string
	createTimezone     string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create an employee",
	Long: `Use this command to create an employee.
"directory:write_users" permission is required (see Y360 help topics).
The password is read from Y360_PASSWORD and is not accepted as a flag.`,
	Example: `  y360c user add --nickname person --department-id 42 --first-name Person --last-name Example`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("user add called")

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

		password, ok := os.LookupEnv("Y360_PASSWORD")
		if !ok || password == "" {
			log.Fatalln("Incorrect settings: password must be specified")
		}

		data, err := helper.NewDirectoryClient(orgId, token).CreateEmployee(helper.CreateEmployeeInput{
			Nickname:     createNickname,
			DepartmentID: createDepartmentID,
			First:        createFirst,
			Last:         createLast,
			Middle:       createMiddle,
			Password:     password,
			About:        createAbout,
			Birthday:     createBirthday,
			DisplayName:  createDisplayName,
			ExternalID:   createExternalID,
			Gender:       createGender,
			Position:     createPosition,
			Language:     createLanguage,
			Timezone:     createTimezone,
		})
		if err != nil {
			log.Fatalln("Unable to create user:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id})
		t.AppendRow(table.Row{"Email", data.Email})
		t.AppendRow(table.Row{"Name", strings.TrimSpace(fmt.Sprintf("%s %s %s", data.Name.Last, data.Name.First, data.Name.Middle))})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	addCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	addCmd.Flags().StringVar(&createNickname, "nickname", "", "employee login")
	addCmd.Flags().IntVar(&createDepartmentID, "department-id", 0, "department id required by the directory create method")
	addCmd.Flags().StringVar(&createFirst, "first-name", "", "employee first name")
	addCmd.Flags().StringVar(&createLast, "last-name", "", "employee last name")
	addCmd.Flags().StringVar(&createMiddle, "middle-name", "", "employee middle name")
	addCmd.Flags().StringVar(&createAbout, "about", "", "about an user")
	addCmd.Flags().StringVar(&createBirthday, "birthday", "", "user birthday (YYYY-MM-DD)")
	addCmd.Flags().StringVar(&createDisplayName, "display-name", "", "user display name")
	addCmd.Flags().StringVar(&createExternalID, "external-id", "", "user external id")
	addCmd.Flags().StringVar(&createGender, "gender", "", "user gender")
	addCmd.Flags().StringVar(&createPosition, "position", "", "user position")
	addCmd.Flags().StringVar(&createLanguage, "language", "", "user language")
	addCmd.Flags().StringVar(&createTimezone, "timezone", "", "user timezone")

	for _, name := range []string{"nickname", "department-id", "first-name", "last-name"} {
		if err := addCmd.MarkFlagRequired(name); err != nil {
			log.Fatalln("Error marking flag as required:", err)
		}
	}
}
