package metadata

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
)

type CargoPackageContent struct {
	Edition string
	Name    string
	Version string
}

type CargoContent struct {
	Pkg          CargoPackageContent `toml:"package,omitempty"`
	Dependencies interface{}
}

func readCargoToml() (*CargoContent, error) {
	var cc CargoContent
	data, err := os.ReadFile("Cargo.toml")
	if err != nil {
		return nil, err
	}

	_, err = toml.Decode(string(data), &cc)
	if err != nil {
		return nil, err
	}

	return &cc, nil
}

func (cfg *Metadata) ExtractProjectNameFromCargoFile(fs *afero.Fs) (*string, error) {
	s, err := readCargoToml()

	if err != nil {
		return nil, err
	}

	return &s.Pkg.Name, nil
}
