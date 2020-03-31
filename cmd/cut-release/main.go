package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/plaid/go-envvar/envvar"
)

var functionDocsTemplate = `
# Functions

## loadMeshStreamingWithURLAsync
▸ **loadMeshStreamingWithURLAsync**(` + "`" + `url` + "`" + `: ` + "`" + `string` + "`" + `): *Promise‹` + "`" + `void` + "`" + `›*

*Defined in [index.ts:7](https://github.com/0xProject/0x-mesh/blob/%s/packages/browser-lite/src/index.ts#L7)*

Loads the Wasm module that is provided by fetching a url.

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
` + "`" + `url` + "`" + ` | ` + "`" + `string` + "`" + ` | The URL to query for the Wasm binary |

<hr />

## loadMeshStreamingAsync

▸ **loadMeshStreamingAsync**(` + "`" + `response` + "`" + `: ` + "`" + `Response | Promise<Response>` + "`" + `): *Promise‹` + "`" + `void` + "`" + `›*

*Defined in [index.ts:15](https://github.com/0xProject/0x-mesh/blob/%s/packages/browser-lite/src/index.ts#L15)*

Loads the Wasm module that is provided by a response.

**Parameters:**

Name | Type | Description |
------ | ------ | ------ |
` + "`" + `response` + "`" + ` | ` + "`" + `Response &#124; Promise<Response>` + "`" + ` | The Wasm response that supplies the Wasm binary |

<hr />`

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

	generateTypescriptDocs()
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

func generateTypescriptDocs() {
	// Generate the initial docs for the Typescript packages. These docs will
	// be used to create the final set of docs.
	cmd := exec.Command("yarn", "docs:md")
	cmd.Dir = "."
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(stdoutStderr))
		log.Fatal(err)
	}
	commitHash, err := getDocsCommitHash("docs/browser-bindings/browser-lite/reference.md")
	if err != nil {
		log.Fatal(err)
	}

	// Copy the browser-lite docs to the `@0x/mesh-browser` packages's `reference.md`
	// file. These docs are the correct docs for the `@0x/mesh-browser` package.
	cmd = exec.Command(
		"cp",
		"docs/browser-bindings/browser-lite/reference.md",
		"docs/browser-bindings/browser/reference.md",
	)
	cmd.Dir = "."
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Print(string(stdoutStderr))
		log.Fatal(err)
	}

	// Create the documentation for the `loadMeshStreamingAsync` and the `loadMeshStreamingWithURLAsync`
	// functions. Append these docs to the end of the existing browser-lite docs.
	functionDocs := fmt.Sprintf(functionDocsTemplate, commitHash, commitHash)
	f, err := os.OpenFile("docs/browser-bindings/browser-lite/reference.md",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(functionDocs); err != nil {
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
	// NOTE(jalextowle): `@0x/mesh-browser` uses the local version of `@0x/mesh-browser-lite`
	// on the `development` branch. Once the `@0x/mesh-browser-lite` package has been published,
	// we need to update dependency in `@0x/mesh-browser` to published version.
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

func getDocsCommitHash(docsPath string) (string, error) {
	dat, err := ioutil.ReadFile(docsPath)
	if err != nil {
		log.Fatal(err)
	}

	regex := "https://github.com/0xProject/0x-mesh/blob/([a-f0-9]+)/"
	var re = regexp.MustCompile(regex)
	matches := re.FindStringSubmatch(string(dat))

	if len(matches) < 2 {
		return "", errors.New("No contents found")
	}
	return matches[1], nil
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
