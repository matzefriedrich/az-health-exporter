package internal

import (
	"fmt"
	"strings"
)

var (
	// Version indicates the current version of the application, typically set at build time or defaulting to an empty string.
	Version = ""

	// CommitSha represents the commit SHA of the current build, defaulting to an empty string if not provided.
	CommitSha = ""

	// ReleaseDate holds the release date information as a string.
	ReleaseDate = ""

	// ReleaseName specifies the name of the release, used to identify the application version or deployment.
	ReleaseName = "az-health-monitor"
)

func GetInformativeApplicationName() string {
	if ReleaseName != "" && Version != "" {
		buffer := strings.Builder{}
		buffer.WriteString(fmt.Sprintf("%s, %s", ReleaseName, Version))
		if CommitSha != "" {
			buffer.WriteString(fmt.Sprintf("-%s", CommitSha))
		}
		if ReleaseDate != "" {
			buffer.WriteString(fmt.Sprintf(" (%s)", ReleaseDate))
		}
		return buffer.String()
	}
	return "az-health-exporter"
}
