package user

import (
	"log"

	"github.com/foxsoft2005/y360c/helper"
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
	createPassword     string
	createOutput       string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create an employee",
	Long: `Use this command to create an employee.
"directory:write_users" permission is required (see Y360 help topics).

Required: --nickname (login), --department-id, --first-name and --last-name.
A random initial password is generated unless --password is set. The employee must change it at the first sign-in.
--output appends nickname, email, user id and password to a CSV file. The same path appends another row.

Optional profile fields, sent only when the flag is present:
  --middle-name     patronymic
  --position        job title
  --display-name    name shown in the profile
  --about           profile description
  --birthday        date as YYYY-MM-DD
  --gender          gender value accepted by the directory
  --language        language code, for example ru
  --timezone        time zone, for example Europe/Moscow
  --external-id     external identifier`,
	Example: `  y360c user add --nickname person --department-id 42 --first-name Person --last-name Example --middle-name Middle --position Engineer --display-name "Person Example" --about "Directory account" --birthday 1990-01-31 --gender unspecified --language ru --timezone Europe/Moscow --external-id ext-100 --output creds.csv
  y360c user add --nickname person --department-id 42 --first-name Person --last-name Example --password TempPass-1`,
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

		password := createPassword
		if password == "" {
			generated, err := helper.GenerateTemporaryPassword()
			if err != nil {
				log.Fatalln("Unable to generate password:", err)
			}
			password = generated
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

		row := credential{
			Nickname: data.Nickname,
			Email:    data.Email,
			UserID:   data.Id,
			Password: password,
		}
		if row.Nickname == "" {
			row.Nickname = createNickname
		}
		renderCredentials([]credential{row})
		if createOutput != "" {
			if err := appendCredentials(createOutput, []credential{row}); err != nil {
				log.Fatalln("Unable to write credentials:", err)
			}
		}
	},
}

func init() {
	addCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	addCmd.Flags().StringVar(&createNickname, "nickname", "", "login")
	addCmd.Flags().IntVar(&createDepartmentID, "department-id", 0, "department id")
	addCmd.Flags().StringVar(&createFirst, "first-name", "", "given name")
	addCmd.Flags().StringVar(&createLast, "last-name", "", "family name")
	addCmd.Flags().StringVar(&createMiddle, "middle-name", "", "patronymic")
	addCmd.Flags().StringVar(&createAbout, "about", "", "profile description")
	addCmd.Flags().StringVar(&createBirthday, "birthday", "", "birthday, YYYY-MM-DD")
	addCmd.Flags().StringVar(&createDisplayName, "display-name", "", "name shown in the profile")
	addCmd.Flags().StringVar(&createExternalID, "external-id", "", "external identifier")
	addCmd.Flags().StringVar(&createGender, "gender", "", "gender value accepted by the directory")
	addCmd.Flags().StringVar(&createPosition, "position", "", "job title")
	addCmd.Flags().StringVar(&createLanguage, "language", "", "language code, for example ru")
	addCmd.Flags().StringVar(&createTimezone, "timezone", "", "time zone, for example Europe/Moscow")
	addCmd.Flags().StringVar(&createPassword, "password", "", "initial password; generated when omitted")
	addCmd.Flags().StringVar(&createOutput, "output", "", "CSV file for created credentials; several accounts are appended as table rows")

	for _, name := range []string{"nickname", "department-id", "first-name", "last-name"} {
		if err := addCmd.MarkFlagRequired(name); err != nil {
			log.Fatalln("Error marking flag as required:", err)
		}
	}
}
