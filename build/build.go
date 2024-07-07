package build

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	// ShortVersion 短版本号
	ShortVersion = "1.7.0"
)

// The value of variables come form `gb build -ldflags '-X "build.Built=xxxxx" -X "build.CommitID=xxxx"' `
var (
	// Built build time
	Built string
	// GitBranch current git branch
	GitBranch string
	// GitCommit git commit id
	GitCommit string
)

// Version 生成版本信息
func Version() string {
	var buf strings.Builder
	buf.WriteString(ShortVersion)

	if Built != "" {
		buf.WriteString(fmt.Sprintf("\n%-15s%s", "Built:", Built))
	}
	if GitBranch != "" {
		buf.WriteString(fmt.Sprintf("\n%-15s%s", "Git branch:", GitBranch))
	}
	if GitCommit != "" {
		buf.WriteString(fmt.Sprintf("\n%-15s%s", "Git commit:", GitCommit))
	}
	buf.WriteString(fmt.Sprintf("\n%-15s%s", "Go version:", runtime.Version()))
	buf.WriteString(fmt.Sprintf("\n%-15s%s/%s", "OS/Arch:", runtime.GOOS, runtime.GOARCH))
	buf.WriteString(fmt.Sprintf("\n%-15s%t", "Experimental:", strings.EqualFold(os.Getenv("G_EXPERIMENTAL"), "true")))
	return buf.String()
}
