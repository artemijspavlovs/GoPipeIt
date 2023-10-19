package metadata

import (
	"errors"
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/afero"
)

type Metadata struct {
	ProgrammingLanguage string
	ProjectName         string
	GoVersion           string
	GitBranch           string
	CICDPlatform        string
	PipelineTasks       map[string]bool
	LocalTasks          map[string]bool
}

func New() *Metadata {
	return &Metadata{}
}

func (cfg *Metadata) SetProgrammingLanguage(n string) {
	cfg.ProgrammingLanguage = n
}

func (cfg *Metadata) SetProjectName(n string, fs *afero.Fs) error {
	switch cfg.ProgrammingLanguage {
	case "go":
		// TODO: extract into a separate function
		err := cfg.SetProjectNameForGoProject(n, fs)
		if err != nil {
			pterm.Error.PrintOnError(err)
			return err
		}

		cfg.ProjectName = n
		return nil
	case "rust":
		fmt.Println("tbd")
	}
	// TODO: improve error message
	return errors.New("something went wrong")
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
