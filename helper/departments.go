package helper

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

// DepartmentRecord is one department from the current directory list method.
type DepartmentRecord struct {
	ID           int
	ParentID     int
	Name         string
	Email        string
	Aliases      []string
	MembersCount int
	Removed      bool
}

type departmentPage struct {
	Total int              `json:"total"`
	Items []departmentItem `json:"items"`
}

type departmentItem struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Aliases      []string `json:"aliases"`
	MembersCount int      `json:"members_count"`
	Removed      bool     `json:"removed"`
	ParentID     int      `json:"parent_id"`
	Parent       *struct {
		ID int `json:"id"`
	} `json:"parent"`
}

// ListDepartments reads every page of the current department list on cloud-api.yandex.net.
func (c DirectoryClient) ListDepartments(limit int) ([]DepartmentRecord, error) {
	if limit <= 0 {
		limit = 100
	}
	var records []DepartmentRecord
	for offset := 0; offset < 100000; offset += limit {
		query := url.Values{}
		query.Set("limit", strconv.Itoa(limit))
		query.Set("offset", strconv.Itoa(offset))
		endpoint := fmt.Sprintf("%s/v1/directory/organizations/%d/departments?%s", c.cloudBase(), c.OrgID, query.Encode())
		resp, err := c.do(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}
		var page departmentPage
		if err := json.Unmarshal(resp.Body, &page); err != nil {
			return nil, err
		}
		if len(page.Items) == 0 {
			break
		}
		for _, item := range page.Items {
			records = append(records, item.record())
		}
		if len(page.Items) < limit || (page.Total > 0 && offset+len(page.Items) >= page.Total) {
			break
		}
	}
	return records, nil
}

func (item departmentItem) record() DepartmentRecord {
	parentID := item.ParentID
	if item.Parent != nil {
		parentID = item.Parent.ID
	}
	return DepartmentRecord{
		ID:           item.ID,
		ParentID:     parentID,
		Name:         item.Name,
		Email:        item.Email,
		Aliases:      item.Aliases,
		MembersCount: item.MembersCount,
		Removed:      item.Removed,
	}
}

// FormatDepartmentTree renders departments as a parent/child list.
// Each line is id, name, mailing address and aliases. Removed departments are omitted.
// Siblings are ordered by id.
func FormatDepartmentTree(items []DepartmentRecord) string {
	return formatDepartmentTree(items, false)
}

// FormatDepartmentTreeByName is FormatDepartmentTree with siblings ordered by name.
func FormatDepartmentTreeByName(items []DepartmentRecord) string {
	return formatDepartmentTree(items, true)
}

func formatDepartmentTree(items []DepartmentRecord, byName bool) string {
	visible := make([]DepartmentRecord, 0, len(items))
	byID := make(map[int]DepartmentRecord, len(items))
	for _, item := range items {
		if item.Removed || item.ID == 0 {
			continue
		}
		visible = append(visible, item)
		byID[item.ID] = item
	}
	children := map[int][]DepartmentRecord{}
	isChild := map[int]bool{}
	for _, item := range visible {
		if item.ParentID == 0 || item.ParentID == item.ID {
			continue
		}
		if _, ok := byID[item.ParentID]; !ok {
			continue
		}
		children[item.ParentID] = append(children[item.ParentID], item)
		isChild[item.ID] = true
	}
	var roots []DepartmentRecord
	for _, item := range visible {
		if !isChild[item.ID] {
			roots = append(roots, item)
		}
	}
	less := func(a, b DepartmentRecord) bool {
		if byName && a.Name != b.Name {
			return a.Name < b.Name
		}
		return a.ID < b.ID
	}
	sort.Slice(roots, func(i, j int) bool { return less(roots[i], roots[j]) })
	for id := range children {
		sort.Slice(children[id], func(i, j int) bool { return less(children[id][i], children[id][j]) })
	}

	var b strings.Builder
	seen := map[int]bool{}
	var walk func(DepartmentRecord, int, map[int]bool)
	walk = func(node DepartmentRecord, depth int, stack map[int]bool) {
		if stack[node.ID] || seen[node.ID] {
			return
		}
		stack[node.ID] = true
		seen[node.ID] = true
		b.WriteString(strings.Repeat("  ", depth))
		b.WriteString(strconv.Itoa(node.ID))
		b.WriteString("  ")
		b.WriteString(node.Name)
		if node.Email != "" {
			b.WriteString("  ")
			b.WriteString(node.Email)
		}
		for _, alias := range node.Aliases {
			if alias == "" {
				continue
			}
			b.WriteString("  ")
			b.WriteString(alias)
		}
		b.WriteByte('\n')
		for _, child := range children[node.ID] {
			walk(child, depth+1, stack)
		}
		delete(stack, node.ID)
	}
	for _, root := range roots {
		walk(root, 0, map[int]bool{})
	}
	return b.String()
}

// DepartmentSubtree returns rootID and its descendants. The selected root is shown at the top.
func DepartmentSubtree(items []DepartmentRecord, rootID int) ([]DepartmentRecord, bool) {
	if rootID == 0 {
		return items, true
	}
	index := make(map[int]DepartmentRecord, len(items))
	children := map[int][]int{}
	for _, item := range items {
		index[item.ID] = item
		children[item.ParentID] = append(children[item.ParentID], item.ID)
	}
	if _, ok := index[rootID]; !ok {
		return nil, false
	}
	var out []DepartmentRecord
	seen := map[int]bool{}
	var walk func(int)
	walk = func(id int) {
		if seen[id] {
			return
		}
		seen[id] = true
		item := index[id]
		if id == rootID {
			item.ParentID = 0
		}
		out = append(out, item)
		for _, child := range children[id] {
			walk(child)
		}
	}
	walk(rootID)
	return out, true
}
