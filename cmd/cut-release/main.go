package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/plaid/go-envvar/envvar"
)

type envVars struct {
	// Version is the new release version to use
	Version string `envvar:"VERSION"`
}

func main() {
	env := envVars{}
	if err := envvar.Parse(&env); err != nil {
		log.Fatal(err)
	}

	updateHardCodedVersions(env.Version)

	// Run `yarn install` to make sure `TypeDoc` dep is installed
	cmd := exec.Command("yarn", "install", "--frozen-lockfile")
	cmd.Dir = "."
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(stdoutStderr))
		log.Fatal(err)
	}

	generateTypescriptClientDocs()
	generateTypescriptBrowserDocs()
	generateTypescriptBrowserLiteDocs()

	createReleaseChangelog(env.Version)
}

func createReleaseChangelog(version string) {
	regex := fmt.Sprintf(`(?ms)(## v%s\n)(.*?)(## v)`, version)
	changelog, err := getFileContentsWithRegex("CHANGELOG.md", regex)
	if err != nil {
		log.Println("No CHANGELOG entries found for version", version)
		return // Noop
	}

	releaseChangelog := fmt.Sprintf(`- [Docker image](https://hub.docker.com/r/0xorg/mesh/tags)
- [README](https://github.com/0xProject/0x-mesh/blob/v%s/README.md)

## Summary
%s
`, version, changelog)

	err = ioutil.WriteFile("RELEASE_CHANGELOG.md", []byte(releaseChangelog), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func generateTypescriptClientDocs() {
	// Run `yarn docs:md` to generate MD docs
	cmd := exec.Command("yarn", "docs:md")
	cmd.Dir = "packages/rpc-client"
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(stdoutStderr))
		log.Fatal(err)
	}
}

func generateTypescriptBrowserDocs() {
	// Run `yarn docs:md` to generate MD docs
	cmd := exec.Command("yarn", "docs:md")
	cmd.Dir = "packages/browser"
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(stdoutStderr))
		log.Fatal(err)
	}
}

func generateTypescriptBrowserLiteDocs() {
	// Run `yarn docs:md` to generate MD docs
	cmd := exec.Command("yarn", "docs:md")
	cmd.Dir = "packages/browser-lite"
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(stdoutStderr))
		log.Fatal(err)
	}
}

// Update the version string in all files that must be updated for a new release
func updateHardCodedVersions(version string) {
	// Update `packages/rpc-client/package.json`
	tsClientPackageJSONPath := "packages/rpc-client/package.json"
	newVersionString := fmt.Sprintf(`"version": "%s"`, version)
	regex := `"version": "(.*)"`
	updateFileWithRegex(tsClientPackageJSONPath, regex, newVersionString)

	// Update `packages/browser-lite/package.json`
	browserLitePackageJSONPath := "packages/browser-lite/package.json"
	newVersionString = fmt.Sprintf(`"version": "%s"`, version)
	regex = `"version": "(.*)"`
	updateFileWithRegex(browserLitePackageJSONPath, regex, newVersionString)

	// Update `packages/browser/package.json`
	browserPackageJSONPath := "packages/browser/package.json"
	newVersionString = fmt.Sprintf(`"version": "%s"`, version)
	regex = `"version": "(.*)"`
	updateFileWithRegex(browserPackageJSONPath, regex, newVersionString)
	newBrowserLiteDependencyString := fmt.Sprintf(`"@0x/mesh-browser-lite": "^%s"`, version)
	regex = `"@0x/mesh-browser-lite": "(.*)"`
	updateFileWithRegex(browserPackageJSONPath, regex, newBrowserLiteDependencyString)

	// Update `core.go`
	corePath := "core/core.go"
	newVersionString = fmt.Sprintf(`version$1= "%s"`, version)
	regex = `version(.*)= "(.*)"`
	updateFileWithRegex(corePath, regex, newVersionString)

	// Update `docs/deployment_with_telemetry.md`
	newVersionString = fmt.Sprintf(`image: 0xorg/mesh:%s`, version)
	regex = `image: 0xorg/mesh:[0-9.]+.*`
	updateFileWithRegex("docs/deployment_with_telemetry.md", regex, newVersionString)

	// Update `CHANGELOG.md`
	changelog := "CHANGELOG.md"
	newChangelogSection := fmt.Sprintf(`## v%s`, version)
	regex = `(## Upcoming release)`
	updateFileWithRegex(changelog, regex, newChangelogSection)

	// Update badge in README.md
	pathToMDFilesWithBadges := []string{"README.md", "docs/rpc_api.md", "docs/deployment.md", "docs/deployment_with_telemetry.md"}
	doubleDashVersion := strings.Replace(version, "-", "--", -1)
	newSvgName := fmt.Sprintf("version-%s-orange.svg", doubleDashVersion)
	regex = `version-(.*)-orange.svg`
	for _, path := range pathToMDFilesWithBadges {
		updateFileWithRegex(path, regex, newSvgName)
	}
}

func updateFileWithRegex(filePath string, regex string, replacement string) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var re = regexp.MustCompile(regex)
	modifiedDat := []byte(re.ReplaceAllString(string(dat), replacement))
	err = ioutil.WriteFile(filePath, modifiedDat, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getFileContentsWithRegex(filePath string, regex string) (string, error) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var re = regexp.MustCompile(regex)
	matches := re.FindAllStringSubmatch(string(dat), -1)

	if len(matches) < 1 || len(matches[0]) < 3 {
		return "", errors.New("No contents found")
	}

	return matches[0][2], nil
}
