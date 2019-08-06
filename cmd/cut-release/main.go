package main

import (
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

	addChangelogUpcomingReleaseTemplate()

	generateTypescriptClientDocs()
}

func addChangelogUpcomingReleaseTemplate() {
	changelogPath := "CHANGELOG.md"
	changelog, err := ioutil.ReadFile(changelogPath)
	if err != nil {
		log.Fatal(err)
	}

	finalChangelogLines := []string{}
	lines := strings.Split(string(changelog), "\n")
	for _, l := range lines {
		finalChangelogLines = append(finalChangelogLines, l)
		if strings.Contains(l, "# CHANGELOG") {
			finalChangelogLines = append(finalChangelogLines, "")
			finalChangelogLines = append(finalChangelogLines, "## Upcoming release")
		}
	}

	updatedChangelog := strings.Join(finalChangelogLines, "\n")
	err = ioutil.WriteFile(changelogPath, []byte(updatedChangelog), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func generateTypescriptClientDocs() {
	// Run `yarn install` to make sure `TypeDoc` dep is installed
	cmd := exec.Command("yarn", "install", "--frozen-lockfile")
	cmd.Dir = "rpc/clients/typescript"
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(stdoutStderr))
		log.Fatal(err)
	}

	// Run `yarn docs:md` to generate MD docs
	cmd = exec.Command("yarn", "docs:md")
	cmd.Dir = "rpc/clients/typescript"
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Print(string(stdoutStderr))
		log.Fatal(err)
	}
}

// Update the version string in all files that must be updated for a new release
func updateHardCodedVersions(version string) {

	// Update `rpc/clients/typescript/package.json`
	tsClientPackageJSONPath := "rpc/clients/typescript/package.json"
	newVersionString := fmt.Sprintf(`"version": "%s"`, version)
	regex := `"version": "(.*)"`
	updateFileWithRegex(tsClientPackageJSONPath, regex, newVersionString)

	// Update `core.go`
	corePath := "core/core.go"
	newVersionString = fmt.Sprintf(`version$1= "%s"`, version)
	regex = `version(.*)= "(.*)"`
	updateFileWithRegex(corePath, regex, newVersionString)

	// Update `CHANGELOG.md`
	changelog := "CHANGELOG.md"
	newChangelogSection := fmt.Sprintf(`## v%s`, version)
	regex = `(## Upcoming release)`
	updateFileWithRegex(changelog, regex, newChangelogSection)

	// Update `beta_telemetry_node/docker-compose.yml`
	dockerComposePath := "examples/beta_telemetry_node/docker-compose.yml"
	newVersionString = fmt.Sprintf(`image: 0xorg/mesh:%s`, version)
	regex = `image: 0xorg/mesh:(.*)`
	updateFileWithRegex(dockerComposePath, regex, newVersionString)

	// Update badge in README.md
	pathToMDFilesWithBadges := []string{"README.md", "docs/USAGE.md", "docs/DEVELOPMENT.md", "docs/DEPLOYMENT.md"}
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
