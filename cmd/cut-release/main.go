package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	generateTypescriptClientDocs()
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

	// Update SUMMARY.md
	tsClientSummaryPath := "docs/json_rpc_clients/typescript/SUMMARY.md"
	dat, err := ioutil.ReadFile(tsClientSummaryPath)
	if err != nil {
		log.Fatal(err)
	}
	// Modify the paths to be prefixed with `json_rpc_clients/typescript`
	modifiedDat := strings.Replace(string(dat), "](", "](json_rpc_clients/typescript/", -1)
	modifiedDat = strings.Replace(modifiedDat, "](json_rpc_clients/typescript/)", "]()", -1)
	finalTsClientSummary := strings.Replace(modifiedDat, "* [", "  * [", -1)

	// Replace the summary content nested under `Typescript client`
	mainSummaryPath := "docs/SUMMARY.md"
	dat, err = ioutil.ReadFile(mainSummaryPath)
	if err != nil {
		log.Fatal(err)
	}
	finalSummaryLines := []string{}
	lines := strings.Split(string(dat), "\n")
	isCutting := false
	for _, l := range lines {
		if strings.Contains(l, "<!-- END TYPEDOC GENERATED SUMMARY -->") {
			isCutting = false
			finalSummaryLines = append(finalSummaryLines, finalTsClientSummary)
		}
		if !isCutting {
			finalSummaryLines = append(finalSummaryLines, l)
		}
		if strings.Contains(l, "<!-- START TYPEDOC GENERATED SUMMARY -->") {
			isCutting = true
		}
	}
	finalSummary := strings.Join(finalSummaryLines, "\n")
	err = ioutil.WriteFile(mainSummaryPath, []byte(finalSummary), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Remove the nested SUMMARY.MD file
	err = os.Remove(tsClientSummaryPath)
	if err != nil {
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
