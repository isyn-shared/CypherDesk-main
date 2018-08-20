package alias

import (
	"go/build"
	"io/ioutil"
	"runtime"
	"strings"
)

// AnalyzePath changes the path depending on the OS
func AnalyzePath(path string) string {
	if isWindows() {
		return windowsify(path)
	}
	return path
}

func isWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

func windowsify(path string) string {
	// -1 means that we need to replace all entries
	return "\\" + strings.Replace(path, "/", "\\", -1)
}

// ReadFile returns text from file
func ReadFile(path string) (string, error) {
	bs, err := ioutil.ReadFile(build.Default.GOPATH + AnalyzePath("src/CypherDesk-main/"+path))
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
