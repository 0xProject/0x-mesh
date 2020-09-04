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

var functionDocsTemplate = "\n# Functions\n\n## loadMeshStreamingForURLAsync\n▸ **loadMeshStreamingWithURLAsync**(`url`: `string`): *Promise‹`void`›*\n\n*Defined in [index.ts:7](https://github.com/0xProject/0x-mesh/blob/%s/packages/mesh-browser-lite/src/index.ts#L7)*\n\nLoads the Wasm module that is provided by fetching a url.\n\n**Parameters:**\n\nName | Type | Description |\n------ | ------ | ------ |\n`url` | `string` | The URL to query for the Wasm binary |\n\n<hr />\n\n## loadMeshStreamingAsync\n\n▸ **loadMeshStreamingAsync**(`response`: `Response | Promise<Response>`): *Promise‹`void`›*\n\n*Defined in [index.ts:15](https://github.com/0xProject/0x-mesh/blob/%s/packages/mesh-browser-lite/src/index.ts#L15)*\n\nLoads the Wasm module that is provided by a response.\n\n**Parameters:**\n\nName | Type | Description |\n------ | ------ | ------ |\n`response` | `Response &#124; Promise<Response>` | The Wasm response that supplies the Wasm binary |\n\n<hr />"

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

	// Run `yarn prettier` to prettify the newly generated docs
	cmd = exec.Command("yarn", "prettier")
	cmd.Dir = "."
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Print(string(stdoutStderr))
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

const (
	captureVersionString                = `"version": "(.*)"`
	captureMeshBrowserVersionString     = `"@0x/mesh-browser": "(.*)"`
	captureMeshBrowserLiteVersionString = `"@0x/mesh-browser-lite": "(.*)"`
)

// Update the version string in all files that must be updated for a new release
func updateHardCodedVersions(version string) {
	newVersionString := fmt.Sprintf(`"version": "%s"`, version)
	newBrowserLiteDependencyString := fmt.Sprintf(`"@0x/mesh-browser-lite": "^%s"`, version)
	newBrowserDependencyString := fmt.Sprintf(`"@0x/mesh-browser": "^%s"`, version)

	// Update `packages/mesh-graphql-client/package.json`
	tsClientPackageJSONPath := "packages/mesh-graphql-client/package.json"
	updateFileWithRegex(tsClientPackageJSONPath, captureVersionString, newVersionString)
	// NOTE(jalextowle): `@0x/mesh-browser` uses the local version of `@0x/mesh-browser-lite`
	// on the `development` branch. Once the `@0x/mesh-browser-lite` package has been published,
	// we need to update dependency in `@0x/mesh-browser` to published version.
	updateFileWithRegex(tsClientPackageJSONPath, captureMeshBrowserLiteVersionString, newBrowserLiteDependencyString)

	// Update `packages/mesh-browser-lite/package.json`
	browserLitePackageJSONPath := "packages/mesh-browser-lite/package.json"
	updateFileWithRegex(browserLitePackageJSONPath, captureVersionString, newVersionString)

	// Update `packages/mesh-browser/package.json`
	browserPackageJSONPath := "packages/mesh-browser/package.json"
	updateFileWithRegex(browserPackageJSONPath, captureVersionString, newVersionString)
	// NOTE(jalextowle): `@0x/mesh-browser` uses the local version of `@0x/mesh-browser-lite`
	// on the `development` branch. Once the `@0x/mesh-browser-lite` package has been published,
	// we need to update dependency in `@0x/mesh-browser` to published version.
	updateFileWithRegex(browserPackageJSONPath, captureMeshBrowserLiteVersionString, newBrowserLiteDependencyString)

	// Update `packages/mesh-webpack-example-lite/package.json`
	webpackExampleLitePackageJSONPath := "packages/mesh-webpack-example-lite/package.json"
	updateFileWithRegex(webpackExampleLitePackageJSONPath, captureMeshBrowserLiteVersionString, newBrowserLiteDependencyString)

	// Update `packages/mesh-webpack-example/package.json`
	webpackExamplePackageJSONPath := "packages/mesh-webpack-example/package.json"
	updateFileWithRegex(webpackExamplePackageJSONPath, captureMeshBrowserVersionString, newBrowserDependencyString)

	// Update `packages/mesh-integration-tests/package.json`
	integrationTestsPackageJSONPath := "packages/mesh-integration-tests/package.json"
	updateFileWithRegex(integrationTestsPackageJSONPath, captureMeshBrowserVersionString, newBrowserDependencyString)

	// Update `packages/mesh-browser-shim/package.json`
	testWasmPackageJSONPath := "packages/mesh-browser-shim/package.json"
	updateFileWithRegex(testWasmPackageJSONPath, captureMeshBrowserLiteVersionString, newBrowserLiteDependencyString)

	// Update `core.go`
	corePath := "core/core.go"
	newVersionString = fmt.Sprintf(`version$1= "%s"`, version)
	updateFileWithRegex(corePath, `version(.*)= "(.*)"`, newVersionString)

	// Update `docs/deployment_with_telemetry.md`
	newVersionString = fmt.Sprintf(`image: 0xorg/mesh:%s`, version)
	updateFileWithRegex("docs/deployment_with_telemetry.md", `image: 0xorg/mesh:[0-9.]+.*`, newVersionString)

	// Update `CHANGELOG.md`
	changelog := "CHANGELOG.md"
	newChangelogSection := fmt.Sprintf(`## v%s`, version)
	updateFileWithRegex(changelog, `(## Upcoming release)`, newChangelogSection)

	// Update badge in README.md
	pathToMDFilesWithBadges := []string{"README.md", "docs/graphql_api.md", "docs/deployment.md", "docs/deployment_with_telemetry.md"}
	doubleDashVersion := strings.Replace(version, "-", "--", -1)
	newSvgName := fmt.Sprintf("version-%s-orange.svg", doubleDashVersion)
	for _, path := range pathToMDFilesWithBadges {
		updateFileWithRegex(path, `version-(.*)-orange.svg`, newSvgName)
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
