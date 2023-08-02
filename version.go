package main

import (
	"strconv"
	"strings"
)

const Version = "0.1.0"

func VersionInfo() (int, int, int) {
	versionInfo := strings.Split(Version, ".")
	major, _ := strconv.Atoi(versionInfo[0])
	minor, _ := strconv.Atoi(versionInfo[1])
	patch, _ := strconv.Atoi(versionInfo[2])
	return major, minor, patch
}
