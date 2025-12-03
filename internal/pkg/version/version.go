package version

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	gitVersion = "dev"
	gitCommit  = ""
	buildDate  = ""
)

type VersionInfo struct {
	GitVersion string `json:"gitVersion" yaml:"gitVersion"`
	GitCommit  string `json:"gitCommit" yaml:"gitCommit"`
	BuildDate  string `json:"buildDate" yaml:"buildDate"`
	GoVersion  string `json:"goVersion" yaml:"goVersion"`
	Compiler   string `json:"compiler" yaml:"compiler"`
	Platform   string `json:"platform" yaml:"platform"`
}

func Get() *VersionInfo {
	return &VersionInfo{
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func IsDevMode() bool {
	return gitVersion == "dev"
}

func IsMrMode() bool {
	return strings.HasPrefix(gitVersion, "mr")
}
