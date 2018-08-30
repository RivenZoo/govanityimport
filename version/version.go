package version

import (
	"fmt"
	"strings"
	"time"
)

var (
	GitVersion = "<UNKNOWN>"
	Ver        = "<UNKNOWN>"
	BuildTime  = "<UNKNOWN>"
	StartTime  = time.Now().Format("2006-01-02 15:04:05")
)

type Version struct {
	Ver        string `json:"version"`
	StartTime  string `json:"start_time"`
	GitVersion string `json:"git_version"`
	BuildTime  string `json:"build_time"`
}

func GetVersion() Version {
	return Version{
		Ver:        Ver,
		StartTime:  StartTime,
		GitVersion: GitVersion,
		BuildTime:  strings.Replace(BuildTime, "T", " ", -1),
	}
}

func GetVersionString(serverName string) string {
	v := GetVersion()
	format := `< %s >: %s
	StartAt: %s
	GitVersion: %s
	BuildTime; %s
	`

	s := fmt.Sprintf(format, serverName, v.Ver, v.StartTime, v.GitVersion, v.BuildTime)
	return s
}

func PrintVersion(serverName string) {
	fmt.Println(GetVersionString(serverName))
}
