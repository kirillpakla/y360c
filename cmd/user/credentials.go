package user

import (
	"encoding/csv"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

// credential is one created employee login. Several rows form the output table.
type credential struct {
	Nickname string
	Email    string
	UserID   string
	Password string
}

func renderCredentials(rows []credential) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"nickname", "email", "user id", "password"})
	for _, row := range rows {
		t.AppendRow(table.Row{row.Nickname, row.Email, row.UserID, row.Password})
	}
	t.Render()
}

// appendCredentials adds rows to a CSV table. The first write creates the header.
// Later user add runs with the same --output path append another row.
func appendCredentials(path string, rows []credential) error {
	if len(rows) == 0 {
		return nil
	}
	info, err := os.Stat(path)
	newFile := os.IsNotExist(err)
	if err != nil && !newFile {
		return err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if newFile || info.Size() == 0 {
		if err := writer.Write([]string{"nickname", "email", "user_id", "password"}); err != nil {
			return err
		}
	}
	for _, row := range rows {
		if err := writer.Write([]string{row.Nickname, row.Email, row.UserID, row.Password}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
