package sync

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	originSeparator = "origin/"
	headSeparator   = "HEAD"
)

func syncGitRepository(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("Processing Git repository: %s\n", path)

	hasChanges, err := hasLocalChanges(ctx, path)
	if err != nil {
		fmt.Printf("  Warning: cannot check local changes: %v\n", err)
	}
	if hasChanges {
		fmt.Printf("  Repository has local changes, skipping auto-sync\n")
		return nil
	}

	branches, err := getRemoteBranches(ctx, path)
	if err != nil {
		return err
	}

	wp := newPool(3, max(len(branches)/2, 5), func(branch string) error {
		cmd := exec.CommandContext(ctx, "git", "fetch", "origin", fmt.Sprintf("%s:%s", branch, branch))
		cmd.Dir = path
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("    Error fetching branch '%s': %v\nOutput: %s\n", branch, err, string(output))
			return err
		}
		fmt.Printf("    Successfully fetched branch '%s'\n", branch)
		return nil
	})

	wp.Start()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, branch := range branches {
			wp.Add(branch)
		}
		wp.CloseJobs()
	}()
	wg.Wait()
	wp.Stop()

	return nil
}

func hasLocalChanges(ctx context.Context, path string) (bool, error) {
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(output) > 0, nil
}

func getRemoteBranches(ctx context.Context, path string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "-r")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := strings.Split(string(output), "\n")
	// Очищаем и фильтруем
	var result []string
	for _, b := range branches {
		b = strings.TrimSpace(b)
		if b != "" && !strings.Contains(b, headSeparator) {
			result = append(result, strings.TrimPrefix(b, originSeparator))
		}
	}
	return result, nil
}
