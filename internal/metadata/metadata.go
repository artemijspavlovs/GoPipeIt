package metadata

import (
	"bufio"
	"fmt"
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

func readGoMod(fs afero.Fs) ([]string, error) {
	exists, _ := afero.Exists(fs, "go.mod")
	if !exists {
		return nil, fmt.Errorf("go.mod does not exist")
	}
	f, err := fs.Open("go.mod")
	if err != nil {
		return nil, err
	}

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)

	var txt []string
	for sc.Scan() {
		txt = append(txt, sc.Text())
	}
	f.Close()
	return txt, nil
}

func (cfg *Metadata) ExtractProjectNameFromGoModFile(fs afero.Fs) (*string, error) {
	txt, err := readGoMod(fs)

	for _, line := range txt {
		if strings.HasPrefix(line, "module ") {
			trimmed := strings.Replace(line, "module ", "", 1)
			t := trimmed[strings.LastIndex(trimmed, "/")+1:]
			return &t, nil
		}
	}
	return nil, fmt.Errorf("failed to read go.mod: %v", err)
}

func (cfg *Metadata) ExtractGoVersionFromGoModFile(fs afero.Fs) (*string, error) {
	txt, err := readGoMod(fs)

	for _, line := range txt {
		if strings.HasPrefix(line, "go ") {
			t := strings.Replace(line, "go ", "", 1)
			return &t, nil
		}
	}
	return nil, fmt.Errorf("failed to read go.mod file: %v", err)
}
