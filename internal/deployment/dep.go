package deployment

import (
	"errors"
	"esdep/internal/config"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func Deploy(entry config.DeployEntry) error {
	script := entry.Script
	if script == "" {
		return errors.New("script is required for deployment")
	}

	if _, err := os.Stat(entry.Path); os.IsNotExist(err) {
		if err := os.MkdirAll(entry.Path, 0755); err != nil {
			return fmt.Errorf("failed to create deployment path: %v", err)
		}
	}

	sshCmd := fmt.Sprintf("ssh -i '%s' -o StrictHostKeyChecking=no", entry.DeployKey)
	gitEnv := append(os.Environ(), fmt.Sprintf("GIT_SSH_COMMAND=%s", sshCmd))

	gitDir := entry.Path
	repoDir := gitDir

	gitPath := fmt.Sprintf("%s/.git", repoDir)
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		cloneCmd := exec.Command("git", "clone", entry.Repo, repoDir)
		cloneCmd.Env = gitEnv
		if out, err := cloneCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to clone repo: %v, output: %s", err, out)
		}
	} else {
		pullCmd := exec.Command("git", "-C", repoDir, "pull")
		pullCmd.Env = gitEnv
		if out, err := pullCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to pull repo: %v, output: %s", err, out)
		}
	}

	cmd := exec.Command("bash", "-c", entry.Script)
	cmd.Dir = entry.Path
	cmd.Env = os.Environ()

	return cmd.Run()
}

func CheckForRemoteUpdates(entry config.DeployEntry) (hasUpdates bool, err error) {
	repoDir := entry.Path
	gitPath := filepath.Join(repoDir, ".git")
	if _, statErr := os.Stat(gitPath); os.IsNotExist(statErr) {
		return true, nil
	}

	env := os.Environ()
	if entry.DeployKey != "" {
		sshCmd := fmt.Sprintf("ssh -i '%s' -o StrictHostKeyChecking=no", entry.DeployKey)
		env = append(env, fmt.Sprintf("GIT_SSH_COMMAND=%s", sshCmd))
	}

	fetchCmd := exec.Command("git", "-C", repoDir, "fetch")
	fetchCmd.Env = env
	if out, err := fetchCmd.CombinedOutput(); err != nil {
		return false, fmt.Errorf("git fetch: %v, output: %s", err, out)
	}

	localOut, err := exec.Command("git", "-C", repoDir, "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("git rev-parse HEAD: %v, output: %s", err, localOut)
	}
	localHash := strings.TrimSpace(string(localOut))

	remoteOut, err := exec.Command("git", "-C", repoDir, "rev-parse", "@{u}").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("git rev-parse @{u}: %v, output: %s", err, remoteOut)
	}
	remoteHash := strings.TrimSpace(string(remoteOut))

	return localHash != remoteHash, nil
}

func RunUpdateChecker(entries []config.DeployEntry, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run first check immediately so user sees output right away
	checkAndDeploy(entries)

	for range ticker.C {
		checkAndDeploy(entries)
	}
}

func checkAndDeploy(entries []config.DeployEntry) {
	for _, entry := range entries {
		hasUpdates, err := CheckForRemoteUpdates(entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Check for updates failed for %s: %v\n", entry.Path, err)
			continue
		}
		if hasUpdates {
			fmt.Printf("New code is available on the remote repository! (%s)\n", entry.Path)
			fmt.Println("Triggering deployment process...")
			if err := Deploy(entry); err != nil {
				fmt.Fprintf(os.Stderr, "Deploy failed for %s: %v\n", entry.Path, err)
			}
		} else {
			fmt.Printf("%s is up to date.\n", entry.Path)
		}
	}
}
