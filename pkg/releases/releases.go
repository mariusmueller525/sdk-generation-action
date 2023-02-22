package releases

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/speakeasy-api/sdk-generation-action/internal/environment"
)

type ReleasesInfo struct {
	ReleaseVersion         string
	OpenAPIDocVersion      string
	SpeakeasyVersion       string
	OpenAPIDocPath         string
	PythonPackagePublished bool
	PythonPackageName      string
	PythonPackageURL       string
	PythonPath             string
	NPMPackagePublished    bool
	NPMPackageName         string
	NPMPackageUrl          string
	TypescriptPath         string
	GoPackagePublished     bool
	GoPackageURL           string
	GoPath                 string
	PHPPackagePublished    bool
	PHPPackageName         string
	PHPPackageURL          string
	PHPPath                string
	JavaPackagePublished   bool
	JavaPackageName        string
	JavaPackageURL         string
	JavaPath               string
}

func (r ReleasesInfo) String() string {
	releasesOutput := []string{}

	if r.NPMPackagePublished {
		releasesOutput = append(releasesOutput, fmt.Sprintf("- [NPM v%s] https://www.npmjs.com/package/%s/v/%s - %s", r.ReleaseVersion, r.NPMPackageName, r.ReleaseVersion, r.TypescriptPath))
	}

	if r.PythonPackagePublished {
		releasesOutput = append(releasesOutput, fmt.Sprintf("- [PyPI v%s] https://pypi.org/project/%s/%s - %s", r.ReleaseVersion, r.PythonPackageName, r.ReleaseVersion, r.PythonPath))
	}

	if r.GoPackagePublished {
		repoPath := os.Getenv("GITHUB_REPOSITORY")
		releasesOutput = append(releasesOutput, fmt.Sprintf("- [Go v%s] https://github.com/%s/releases/tag/v%s - %s", r.ReleaseVersion, repoPath, r.ReleaseVersion, r.GoPath))
	}

	if r.PHPPackagePublished {
		releasesOutput = append(releasesOutput, fmt.Sprintf("- [Composer v%s] https://packagist.org/packages/%s#v%s - %s", r.ReleaseVersion, r.PHPPackageName, r.ReleaseVersion, r.PHPPath))
	}

	if r.JavaPackagePublished {
		lastDotIndex := strings.LastIndex(r.JavaPackageName, ".")
		groupID := r.JavaPackageName[:lastDotIndex]    // everything before last occurrence of '.'
		artifact := r.JavaPackageName[lastDotIndex+1:] // everything after last occurrence of '.'
		releasesOutput = append(releasesOutput, fmt.Sprintf("- [Maven Central v%s] https://central.sonatype.com/artifact/%s/%s/%s - %s", r.ReleaseVersion, groupID, artifact, r.ReleaseVersion, r.JavaPath))
	}

	if len(releasesOutput) > 0 {
		releasesOutput = append([]string{"\n### Releases"}, releasesOutput...)
	}

	return fmt.Sprintf(`%s## Version %s
### Changes
Based on:
- OpenAPI Doc %s %s
- Speakeasy CLI %s https://github.com/speakeasy-api/speakeasy%s`, "\n\n", r.ReleaseVersion, r.OpenAPIDocVersion, r.OpenAPIDocPath, r.SpeakeasyVersion, strings.Join(releasesOutput, "\n"))
}

func UpdateReleasesFile(releaseInfo ReleasesInfo) error {
	releasesPath := getReleasesPath()

	f, err := os.OpenFile(releasesPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf("error opening releases file: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(releaseInfo.String())
	if err != nil {
		return fmt.Errorf("error writing to releases file: %w", err)
	}

	return nil
}

var (
	releaseInfoRegex     = regexp.MustCompile(`(?s)## Version (\d+\.\d+\.\d+)\n### Changes\nBased on:\n- OpenAPI Doc (.*?) (.*?)\n- Speakeasy CLI (.*?) .*?`)
	npmReleaseRegex      = regexp.MustCompile(`- \[NPM v\d+\.\d+\.\d+\] (https:\/\/www\.npmjs\.com\/package\/(.*?)\/v\/\d+\.\d+\.\d+) - (.*)`)
	pypiReleaseRegex     = regexp.MustCompile(`- \[PyPI v\d+\.\d+\.\d+\] (https:\/\/pypi\.org\/project\/(.*?)\/\d+\.\d+\.\d+) - (.*)`)
	goReleaseRegex       = regexp.MustCompile(`- \[Go v\d+\.\d+\.\d+\] (.*?) - (.*)`)
	composerReleaseRegex = regexp.MustCompile(`- \[Composer v\d+\.\d+\.\d+\] (https:\/\/packagist\.org\/packages\/(.*?)#v\d+\.\d+\.\d+) - (.*)`)
	mavenReleaseRegex    = regexp.MustCompile(`- \[Maven Central v\d+\.\d+\.\d+\] (https:\/\/central\.sonatype\.com\/artifact\/(.*?)\/(.*?)\/.*?) - (.*)`)
)

func GetLastReleaseInfo() (*ReleasesInfo, error) {
	releasesPath := getReleasesPath()

	data, err := os.ReadFile(releasesPath)
	if err != nil {
		return nil, fmt.Errorf("error reading releases file: %w", err)
	}

	return ParseReleases(string(data))
}

func ParseReleases(data string) (*ReleasesInfo, error) {
	releases := strings.Split(data, "\n\n")

	lastRelease := releases[len(releases)-1]

	matches := releaseInfoRegex.FindStringSubmatch(lastRelease)

	fmt.Println(matches)

	if len(matches) != 5 {
		return nil, fmt.Errorf("error parsing last release info")
	}

	info := &ReleasesInfo{
		ReleaseVersion:    matches[1],
		OpenAPIDocVersion: matches[2],
		OpenAPIDocPath:    matches[3],
		SpeakeasyVersion:  matches[4],
	}

	npmMatches := npmReleaseRegex.FindStringSubmatch(lastRelease)

	if len(npmMatches) == 4 {
		info.NPMPackagePublished = true
		info.NPMPackageUrl = npmMatches[1]
		info.NPMPackageName = npmMatches[2]
		info.TypescriptPath = npmMatches[3]
	}

	pypiMatches := pypiReleaseRegex.FindStringSubmatch(lastRelease)

	if len(pypiMatches) == 4 {
		info.PythonPackagePublished = true
		info.PythonPackageURL = pypiMatches[1]
		info.PythonPackageName = pypiMatches[2]
		info.PythonPath = pypiMatches[3]
	}

	goMatches := goReleaseRegex.FindStringSubmatch(lastRelease)

	if len(goMatches) == 3 {
		info.GoPackagePublished = true
		info.GoPackageURL = goMatches[1]
		info.GoPath = goMatches[2]
	}

	composerMatches := composerReleaseRegex.FindStringSubmatch(lastRelease)

	if len(composerMatches) == 4 {
		info.PHPPackagePublished = true
		info.PHPPackageURL = composerMatches[1]
		info.PHPPackageName = composerMatches[2]
		info.PHPPath = composerMatches[3]
	}

	mavenMatches := mavenReleaseRegex.FindStringSubmatch(lastRelease)

	if len(mavenMatches) == 5 {
		info.JavaPackagePublished = true
		info.JavaPackageURL = mavenMatches[1]
		groupID := mavenMatches[2]
		artifact := mavenMatches[3]
		info.JavaPath = mavenMatches[4]

		info.JavaPackageName = fmt.Sprintf(`%s.%s`, groupID, artifact)
	}

	return info, nil
}

func getReleasesPath() string {
	baseDir := environment.GetBaseDir()

	return path.Join(baseDir, "repo", "RELEASES.md")
}
