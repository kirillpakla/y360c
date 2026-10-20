package helper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFormatDepartmentTreeNestsChildrenAndAliases(t *testing.T) {
	got := FormatDepartmentTree([]DepartmentRecord{
		{ID: 21, ParentID: 20, Name: "Ops"},
		{ID: 20, ParentID: 10, Name: "IT", Email: "it@example.com", Aliases: []string{"ops@example.com"}},
		{ID: 10, ParentID: 0, Name: "School", Email: "school@example.com"},
		{ID: 30, ParentID: 10, Name: "HR", Email: "hr@example.com"},
	})
	want := "" +
		"10  School  school@example.com\n" +
		"  20  IT  it@example.com  ops@example.com\n" +
		"    21  Ops\n" +
		"  30  HR  hr@example.com\n"
	if got != want {
		t.Fatalf("tree =\n%s\nwant =\n%s", got, want)
	}
}

func TestFormatDepartmentTreeByNameOrdersSiblings(t *testing.T) {
	got := FormatDepartmentTreeByName([]DepartmentRecord{
		{ID: 20, ParentID: 10, Name: "IT"},
		{ID: 10, ParentID: 0, Name: "School"},
		{ID: 30, ParentID: 10, Name: "HR"},
	})
	want := "" +
		"10  School\n" +
		"  30  HR\n" +
		"  20  IT\n"
	if got != want {
		t.Fatalf("tree =\n%s", got)
	}
}

func TestFormatDepartmentTreeKeepsMissingParentVisible(t *testing.T) {
	got := FormatDepartmentTree([]DepartmentRecord{
		{ID: 8, ParentID: 99, Name: "Detached", Email: "detached@example.com"},
	})
	if !strings.Contains(got, "8  Detached  detached@example.com") {
		t.Fatalf("tree = %q", got)
	}
}

func TestListDepartmentsPagesCloudDirectory(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s", r.Method)
		}
		paths = append(paths, r.URL.RequestURI())
		switch r.URL.Query().Get("offset") {
		case "", "0":
			_, _ = ioWrite(w, `{"limit":1,"offset":0,"total":2,"items":[{"id":10,"name":"School","email":"school@example.com","aliases":["all@example.com"],"members_count":2,"removed":false,"parent":{"id":0}}]}`)
		default:
			_, _ = ioWrite(w, `{"limit":1,"offset":1,"total":2,"items":[{"id":20,"name":"IT","email":"it@example.com","members_count":1,"removed":false,"parent_id":10}]}`)
		}
	}))
	defer server.Close()

	items, err := (DirectoryClient{
		CloudBase: server.URL,
		Token:     "synthetic-token",
		OrgID:     1001,
		HTTP:      server.Client(),
	}).ListDepartments(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].ID != 10 || items[0].ParentID != 0 || items[0].Aliases[0] != "all@example.com" || items[1].ParentID != 10 {
		t.Fatalf("items = %+v", items)
	}
	if len(paths) != 2 || paths[0] != "/v1/directory/organizations/1001/departments?limit=1&offset=0" || paths[1] != "/v1/directory/organizations/1001/departments?limit=1&offset=1" {
		t.Fatalf("paths = %#v", paths)
	}
}

func ioWrite(w http.ResponseWriter, body string) (int, error) {
	return w.Write([]byte(body))
}
