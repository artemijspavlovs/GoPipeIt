package metadata

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/afero"
)

type Metadata struct {
	ProjectName   string
	GoVersion     string
	GitBranch     string
	CICDPlatform  string
	PipelineTasks map[string]bool
	LocalTasks    map[string]bool
}

func New() *Metadata {
	return &Metadata{}
}

func (cfg *Metadata) SetProjectName(n string, fs afero.Fs) error {
	if n == "" {
		fmt.Println("Project name was not set, extracting from go.mod file")
		pn, err := cfg.ExtractProjectNameFromGoModFile(fs)
		if err != nil {
			return fmt.Errorf("failed to extract project name from go.mod file: %v", err)
		}
		pterm.Println("Project name extracted from go.mod:", pterm.Yellow(*pn))
		err = cfg.SetProjectName(*pn, fs)
		if err != nil {
			return fmt.Errorf("failed to set project name: %v", err)
		}
		return nil
	}
	cfg.ProjectName = n
	return nil
}

func (cfg *Metadata) SetGoVersion(n string, fs afero.Fs) error {
	if n == "" {
		fmt.Println("Go version was not set, extracting from go.mod file")
		gv, err := cfg.ExtractGoVersionFromGoModFile(fs)
		if err != nil {
			return fmt.Errorf("failed to extract go version from go.mod file: %v", err)
		}
		pterm.Println("Go version extracted from go.mod: " + pterm.Yellow(*gv))
		err = cfg.SetGoVersion(*gv, fs)
		if err != nil {
			return fmt.Errorf("failed to set go version: %v", err)
		}
		return nil
	}
	cfg.GoVersion = n
	return nil
}

func (cfg *Metadata) SetGitBranch(n string) {
	if n == "" {
		pterm.Println(pterm.White("Git branch was not set, defaulting to ") + pterm.Yellow("main"))
		cfg.SetGitBranch("main")
		return
	}
	cfg.GitBranch = n
}

func (cfg *Metadata) SetCICDPlatform(n string) {
	cfg.CICDPlatform = n
}

func (cfg *Metadata) SetPipelineTasks(n []string) {
	cfg.PipelineTasks = make(map[string]bool)
	for _, t := range n {
		cfg.PipelineTasks[t] = true
	}
}

func (cfg *Metadata) SetLocalTasks(n []string) {
	cfg.LocalTasks = make(map[string]bool)
	for _, t := range n {
		cfg.LocalTasks[t] = true
	}
}

func readGoMod(fs afero.Fs) (*bufio.Scanner, error) {
	exists, _ := afero.Exists(fs, "go.mod")
	if !exists {
		return nil, fmt.Errorf("go.mod does not exist")
	}
	f, err := fs.Open("go.mod")
	if err != nil {
		return nil, err
	}

	sc := bufio.NewScanner(f)
	return sc, nil
}

func (cfg *Metadata) ExtractProjectNameFromGoModFile(fs afero.Fs) (*string, error) {
	sc, err := readGoMod(fs)
	linePrefix := "module "
	var modulePath string

	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, linePrefix) {
			modulePath = strings.TrimSpace(strings.TrimPrefix(line, linePrefix))
			break
		}
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	if modulePath == "" {
		return nil, err
	}

	parts := strings.Split(modulePath, "/")
	lastPart := parts[len(parts)-1]
	if strings.HasPrefix(lastPart, "v") && len(lastPart) > 1 {
		parts = parts[:len(parts)-1]
	}

	moduleName := filepath.Base(strings.Join(parts, "/"))

	return &moduleName, nil
}

func (cfg *Metadata) ExtractGoVersionFromGoModFile(fs afero.Fs) (*string, error) {
	sc, err := readGoMod(fs)
	linePrefix := "go "
	var goVersion string

	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, linePrefix) {
			goVersion = strings.TrimSpace(strings.TrimPrefix(line, linePrefix))
			break
		}
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	if goVersion == "" {
		return nil, err
	}

	// for _, line := range txt {
	// 	if strings.HasPrefix(line, "go ") {
	// 		t := strings.Replace(line, "go ", "", 1)
	// 		return &t, nil
	// 	}
	// }
	return &goVersion, nil
}
