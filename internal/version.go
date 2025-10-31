package internal

import (
	"fmt"
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
	name := "az-health-exporter"
	if ReleaseName != "" && Version != "" {
		name = fmt.Sprintf("%s, %s-%s (%s)", ReleaseName, Version, CommitSha, ReleaseDate)
	}
	return name
}
