package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark-burns-0/devsync/internal/config"
)

type Scanner struct {
	cfg      *config.SyncConfig
	skipDirs map[string]bool
}

func New(cfg *config.SyncConfig) *Scanner {
	return &Scanner{
		cfg: cfg,
		skipDirs: map[string]bool{
			"node_modules":            true,
			"vendor":                  true,
			".cache":                  true,
			"__pycache__":             true,
			"dist":                    true,
			"build":                   true,
			".idea":                   true,
			".vscode":                 true,
			"docker-composer-laravel": true,
			"flutter":                 true,
		},
	}
}

// ScanDirs сканирует директории и возвращает пути к git репозиториям (BFS версия)
func (s *Scanner) ScanDirs() ([]string, error) {
	var dirs []string
	queue := []string{s.cfg.ProjectsRoot}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		entries, err := os.ReadDir(current)
		if err != nil {
			fmt.Printf("Warning: cannot read directory %s: %v\n", current, err)
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			name := entry.Name()
			if s.skipDirs[name] {
				continue
			}

			fullPath := filepath.Join(current, name)
			if name == s.cfg.GitDir {
				projectPath := strings.TrimSuffix(fullPath, s.cfg.GitDir)
				dirs = append(dirs, projectPath)
				continue
			}

			queue = append(queue, fullPath)
		}
	}

	return dirs, nil
}

// ScanDirsDFS сканирует директории и возвращает пути к git репозиториям (DFS версия)
func (s *Scanner) ScanDirsDFS() ([]string, error) {
	dirs := []string{}

	err := filepath.WalkDir(s.cfg.ProjectsRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == s.cfg.GitDir {
			path := strings.TrimSuffix(path, s.cfg.GitDir)
			cleanPath := strings.Trim(path, string(filepath.Separator))
			relPath, _ := filepath.Rel(s.cfg.ProjectsRoot, cleanPath)
			if relPath == "." {
				relPath = ""
			}
			dirs = append(dirs, relPath)
			return fs.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return dirs, nil
}
