package templates

import (
	"os"
	"text/template"

	"github.com/pterm/pterm"
	"github.com/spf13/afero"

	"github.com/artemijspavlovs/gopipeit/internal/metadata"
	"github.com/artemijspavlovs/gopipeit/internal/state"
)

type Templates struct {
	Pairs []state.SourceToDest
}

func New() *Templates {
	return &Templates{}
}

func WriteToFile(t *template.Template, f afero.File, cfg *metadata.Metadata) error {
	err := t.Execute(f, cfg)
	if err != nil {
		return err
	}
	return nil
}

func (t *Templates) AddPairs(pairs []state.SourceToDest) {
	t.Pairs = append(t.Pairs, pairs...)
}

func (t *Templates) CreateDirectoryStructure(p []string, fs afero.Fs) error {
	for _, path := range p {
		pterm.Info.Printfln("creating %s directory", pterm.Yellow(path))
		err := fs.MkdirAll(path, os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	return nil
}
