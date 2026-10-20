// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package dept

import (
	"errors"
	"fmt"
	"log"

	"github.com/goccy/go-json"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/spf13/cobra"
)

type sortingMethod string

const byId sortingMethod = "id"

func (e *sortingMethod) String() string {
	return string(*e)
}

func (e *sortingMethod) Set(v string) error {
	switch v {
	case "id", "name":
		*e = sortingMethod(v)
		return nil
	default:
		return errors.New(`must be one of "id", or "name"`)
	}
}

func (e *sortingMethod) Type() string {
	return "sortingMethod"
}

// Command flags & parameters

var (
	maxRec  int
	orderBy sortingMethod
	asRaw   bool
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Get a hierarchical list of departments",
	Long: `Use this command to retrieve departments of the selected organization as a tree.
Mailing aliases are printed on the same line as the department.
"directory:read_departments" permission is required (see Y360 help topics).`,
	Example: `  y360c dept ls
  y360c dept ls --parent-id 10
  y360c dept ls --order-by name`,
	Run: func(cmd *cobra.Command, args []string) {
		if !asRaw {
			log.Print("dept ls called")
		}

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

		records, err := helper.NewDirectoryClient(orgId, token).ListDepartments(maxRec)
		if err != nil {
			log.Fatalln("Unable to list departments:", err)
		}
		if parentId > 0 {
			selected, ok := helper.DepartmentSubtree(records, parentId)
			if !ok {
				log.Fatalln("Department was not found")
			}
			records = selected
		}

		if asRaw {
			buff, _ := json.MarshalIndent(records, "", "     ")
			fmt.Print(string(buff))
		} else if string(orderBy) == "name" {
			fmt.Print(helper.FormatDepartmentTreeByName(records))
		} else {
			fmt.Print(helper.FormatDepartmentTree(records))
		}

	},
}

func init() {
	orderBy = byId

	lsCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	lsCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	lsCmd.Flags().IntVar(&maxRec, "max", 100, "max records to retrieve")
	lsCmd.Flags().IntVar(&parentId, "parent-id", 0, "parent depratment id")
	lsCmd.Flags().Var(&orderBy, "order-by", "sort by (id or name)")
	lsCmd.Flags().BoolVar(&asRaw, "raw", false, "export as raw data")
}
