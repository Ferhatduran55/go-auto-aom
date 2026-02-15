package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// UpdateInfo güncelleme bilgilerini içerir
type UpdateInfo struct {
	Available      bool   `json:"available"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	ReleaseURL     string `json:"release_url"`
	ReleaseName    string `json:"release_name"`
	PublishedAt    string `json:"published_at"`
	ReleaseNotes   string `json:"release_notes"`
	Error          string `json:"error,omitempty"`
}

// GitHubRelease GitHub API release yanıtı
type GitHubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
	Prerelease  bool   `json:"prerelease"`
	Draft       bool   `json:"draft"`
}

const (
	githubRepo    = "Ferhatduran55/go-auto-aom"
	githubAPIBase = "https://api.github.com/repos/"
)

// checkForUpdates GitHub'dan en son release'i kontrol eder
func checkForUpdates() string {
	info := UpdateInfo{
		CurrentVersion: AppVersion,
	}

	// GitHub API'den latest release'i al
	url := fmt.Sprintf("%s%s/releases/latest", githubAPIBase, githubRepo)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		info.Error = fmt.Sprintf("Request creation failed: %v", err)
		return toJSON(info)
	}

	req.Header.Set("User-Agent", "AutoManagement/"+AppVersion)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		info.Error = fmt.Sprintf("Connection error: %v", err)
		return toJSON(info)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		info.Error = fmt.Sprintf("GitHub API error: %d", resp.StatusCode)
		return toJSON(info)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		info.Error = fmt.Sprintf("Response read error: %v", err)
		return toJSON(info)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		info.Error = fmt.Sprintf("JSON parse error: %v", err)
		return toJSON(info)
	}

	// Versiyon karşılaştır
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	info.LatestVersion = latestVersion
	info.ReleaseURL = release.HTMLURL
	info.ReleaseName = release.Name
	info.PublishedAt = release.PublishedAt
	info.ReleaseNotes = release.Body

	// Semantic versioning karşılaştırması
	if isNewerVersion(latestVersion, AppVersion) {
		info.Available = true
	}

	return toJSON(info)
}

// isNewerVersion a > b mi kontrol eder (semantic versioning)
func isNewerVersion(latest, current string) bool {
	latestParts := parseVersion(latest)
	currentParts := parseVersion(current)

	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true
		}
		if latestParts[i] < currentParts[i] {
			return false
		}
	}
	return false
}

// parseVersion "25.12.3" -> [25, 12, 3]
func parseVersion(v string) [3]int {
	parts := strings.Split(v, ".")
	var result [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		result[i], _ = strconv.Atoi(parts[i])
	}
	return result
}

// toJSON helper
func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// getAppVersion uygulama versiyonunu döndürür
func getAppVersion() string {
	return AppVersion
}
