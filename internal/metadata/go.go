package metadata

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"golang.org/x/mod/modfile"
)

func readGoMod() (*modfile.File, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return nil, err
	}

	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (cfg *Metadata) ExtractProjectNameFromGoModFile(fs *afero.Fs) (*string, error) {
	f, err := readGoMod()
	if err != nil {
		return nil, err
	}
	name := f.Module.Mod.Path[strings.LastIndex(f.Module.Mod.Path, "/")+1:]

	return &name, nil
}

func (cfg *Metadata) ExtractGoVersionFromGoModFile(fs *afero.Fs) (*string, error) {
	f, err := readGoMod()
	if err != nil {
		return nil, err
	}

	return &f.Go.Version, nil
}

func (cfg *Metadata) SetGoVersion(n string, fs *afero.Fs) error {
	if n == "" {
		fmt.Println("Go version was not set, extracting from go.mod file")
		gv, err := cfg.ExtractGoVersionFromGoModFile(fs)
		if err != nil {
			return fmt.Errorf("failed to extract go version from go.mod file: %v", err)
		}
		pterm.Println("Go version extracted from go.mod: " + pterm.Yellow(*gv))
		cfg.GoVersion = *gv
		return nil
	}
	cfg.GoVersion = n
	return errors.New("something went wrong")
}

func (cfg *Metadata) SetProjectNameForGoProject(n string, fs *afero.Fs) error {
	if n == "" {
		fmt.Println("Project name was not set, extracting from go.mod file")
		pn, err := cfg.ExtractProjectNameFromGoModFile(fs)
		if err != nil {
			return fmt.Errorf("failed to extract project name from go.mod file: %v", err)
		}
		pterm.Println("Project name extracted from go.mod:", pterm.Yellow(*pn))
		cfg.ProjectName = *pn
		return nil
	}
	cfg.GoVersion = n
	return errors.New("something went wrong")
}
