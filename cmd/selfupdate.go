package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/rkriad585/neovector/internal/config"
	"github.com/rkriad585/neovector/internal/version"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update neovector to the latest version",
	Long: `Fetches the latest version from GitHub and updates the neovector binary
if a newer release is available.

The update process is safe and platform-aware:
  - Windows: renames the running binary to .old before replacing
  - Unix: replaces the binary in-place`,
	RunE: selfUpdateRun,
}

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
	selfUpdateCmd.Flags().StringP("proxy", "p", "", "Proxy URL for the self-update (e.g. http://proxy:8080)")
}

func selfUpdateRun(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Println("Checking for neovector updates...")
	fmt.Println()

	proxyURL := resolveProxy(cmd)
	if proxyURL != "" {
		Yellow.Printf("  Using proxy: %s\n", proxyURL)
		os.Setenv("HTTP_PROXY", proxyURL)
		os.Setenv("HTTPS_PROXY", proxyURL)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	remoteVersion, err := fetchRemoteVersion(client, cmd)
	if err != nil {
		return err
	}

	Cyan.Printf("  Current version: %s\n", version.Version)
	Green.Printf("  Latest version : %s\n", remoteVersion)
	fmt.Println()

	if strings.TrimSpace(remoteVersion) == strings.TrimSpace(version.Version) {
		Green.Println("  neovector is already up to date!")
		return nil
	}

	Yellow.Printf("  Updating neovector from %s to %s...\n", version.Version, remoteVersion)

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to locate running executable: %w", err)
	}

	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	fileName := fmt.Sprintf("%s-%s-%s%s", config.ProjectName, runtime.GOOS, runtime.GOARCH, ext)
	downloadURL := fmt.Sprintf("https://github.com/rkriad585/%s/releases/download/%s/%s",
		config.ProjectName, remoteVersion, fileName)

	Cyan.Printf("  Downloading from: %s\n", downloadURL)

	dlClient := &http.Client{Timeout: 5 * time.Minute}
	dlReq, err := http.NewRequestWithContext(cmd.Context(), "GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}
	dlResp, err := dlClient.Do(dlReq)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer dlResp.Body.Close()

	if dlResp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download release binary: server returned %s", dlResp.Status)
	}

	if err := replaceBinary(exePath, remoteVersion, dlResp.Body); err != nil {
		return err
	}

	return nil
}

func resolveProxy(cmd *cobra.Command) string {
	if p, _ := cmd.Flags().GetString("proxy"); p != "" {
		return p
	}
	if appCfg != nil && appCfg.Network.Proxy != "" {
		return appCfg.Network.Proxy
	}
	return ""
}

func fetchRemoteVersion(client *http.Client, cmd *cobra.Command) (string, error) {
	versionURL := fmt.Sprintf("https://raw.githubusercontent.com/rkriad585/%s/main/.version",
		config.ProjectName)
	req, err := http.NewRequestWithContext(cmd.Context(), "GET", versionURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to check remote version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to check remote version: server returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read version info: %w", err)
	}

	remoteVersion := strings.TrimSpace(string(body))
	if remoteVersion == "" {
		return "", fmt.Errorf("received empty version info from server")
	}
	return remoteVersion, nil
}

func replaceBinary(exePath, remoteVersion string, body io.Reader) error {
	tempPath := exePath + ".tmp"
	os.Remove(tempPath)

	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0755)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	cleanup := true
	defer func() {
		if cleanup {
			tempFile.Close()
			os.Remove(tempPath)
		}
	}()

	progressWriter := &WriteCounter{}
	teeReader := io.TeeReader(body, progressWriter)

	if _, err := io.Copy(tempFile, teeReader); err != nil {
		return fmt.Errorf("failed to write binary content: %w", err)
	}
	tempFile.Close()
	fmt.Println()

	oldPath := exePath + ".old"
	os.Remove(oldPath)

	if runtime.GOOS == "windows" {
		if err := os.Rename(exePath, oldPath); err != nil {
			return fmt.Errorf("failed to rename running executable: %w", err)
		}
		if err := os.Rename(tempPath, exePath); err != nil {
			os.Rename(oldPath, exePath)
			return fmt.Errorf("failed to install new executable: %w", err)
		}
		cleanup = false
		Green.Printf("  Success! neovector has been updated to %s.\n", remoteVersion)
		Yellow.Println("  Note: You can safely delete", oldPath, "after closing this session.")
	} else {
		if err := os.Rename(tempPath, exePath); err != nil {
			return fmt.Errorf("failed to install new executable: %w", err)
		}
		cleanup = false
		os.Chmod(exePath, 0755)
		Green.Printf("  Success! neovector has been updated to %s.\n", remoteVersion)
	}

	return nil
}

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	fmt.Printf("\r  Downloaded: %.2f MB", float64(wc.Total)/1024/1024)
	return n, nil
}
