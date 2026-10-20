package helper

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// ErrTokenMissing is returned when Y360_TOKEN is absent from the process environment.
// The error text never includes the token or the organization id.
var ErrTokenMissing = errors.New("access token must be specified")

// Settings holds non-secret directory settings and the process token.
// Token is filled only from the environment lookup, never from the file.
type Settings struct {
	OrgID   int
	APIBase string
	Token   string
}

var (
	directoryBase = BaseUrl
	cloudBase     = BaseUrlV2
)

// LoadSettings reads Y360_ORG_ID and Y360_API_BASE from envFile.
// Y360_TOKEN in the file is ignored. lookup supplies the process environment.
func LoadSettings(envFile string, lookup func(string) (string, bool)) (Settings, error) {
	values, err := parseEnvFile(envFile)
	if err != nil {
		return Settings{}, err
	}

	var settings Settings
	if raw := strings.TrimSpace(values["Y360_ORG_ID"]); raw != "" {
		id, convErr := strconv.Atoi(raw)
		if convErr != nil || id < 0 {
			return Settings{}, errors.New("organization id is invalid")
		}
		settings.OrgID = id
	}
	settings.APIBase = strings.TrimSpace(values["Y360_API_BASE"])

	if lookup != nil {
		if token, ok := lookup("Y360_TOKEN"); ok {
			settings.Token = token
		}
	}
	if settings.Token == "" {
		return settings, ErrTokenMissing
	}
	return settings, nil
}

func parseEnvFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, errors.New("settings file contains a line without a separator")
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(value) >= 2 {
			quote := value[0]
			if (quote == '"' || quote == '\'') && value[len(value)-1] == quote {
				value = value[1 : len(value)-1]
			}
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}

// OrgIDStatus reports whether an organization id is present without revealing it.
func OrgIDStatus(id int) string {
	if id == 0 {
		return "unset"
	}
	return "set"
}

// ApplySettings publishes non-secret settings for commands.
// A non-empty API base replaces both documented hosts for this process so tests stay local.
// An empty API base restores the documented host of each method.
func ApplySettings(settings Settings) {
	viper.Set("org-id", settings.OrgID)
	if settings.APIBase != "" {
		base := strings.TrimRight(settings.APIBase, "/")
		directoryBase = base
		cloudBase = base
		viper.Set("api-base", base)
		return
	}
	directoryBase = BaseUrl
	cloudBase = BaseUrlV2
	viper.Set("api-base", "")
}

// DirectoryBase is the host for Directory API v1 methods on api360.yandex.net.
func DirectoryBase() string {
	return directoryBase
}

// CloudBase is the host for methods documented on cloud-api.yandex.net.
func CloudBase() string {
	return cloudBase
}
