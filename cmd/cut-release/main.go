package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
		log.Print(stdoutStderr)
		log.Fatal(err)
	}

	// Run `yarn docs:md` to generate MD docs
	cmd = exec.Command("yarn", "docs:md")
	cmd.Dir = "rpc/clients/typescript"
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Print(stdoutStderr)
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
		if strings.Contains(l, "## Contributing") {
			isCutting = false
			finalSummaryLines = append(finalSummaryLines, finalTsClientSummary, "")
		}
		if !isCutting {
			finalSummaryLines = append(finalSummaryLines, l)
		}
		if strings.Contains(l, "[Typescript client]") {
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
	updateVersionKey(tsClientPackageJSONPath, version)

	// Update `core.ts`
	corePath := "core/core.go"
	updateVersionKey(corePath, version)

	// Update badge in README.md
	readmePath := "README.md"
	updateBadge(readmePath, version)
}

func updateVersionKey(filePath string, version string) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	newVersionString := fmt.Sprintf(`"version": "%s",`, version)
	modifiedDat := []byte(strings.Replace(string(dat), `"version": "development",`, newVersionString, 1))
	err = ioutil.WriteFile(filePath, modifiedDat, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func updateBadge(filePath string, version string) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	doubleDashVersion := strings.Replace(version, "-", "--", -1)
	newSvgName := fmt.Sprintf("version-%s-orange.svg", doubleDashVersion)
	modifiedDat := []byte(strings.Replace(string(dat), "version-development-orange.svg", newSvgName, 1))
	err = ioutil.WriteFile(filePath, modifiedDat, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
