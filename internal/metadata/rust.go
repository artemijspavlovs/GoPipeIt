package metadata

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
)

type CargoPackageContent struct {
	Edition        string `toml:"edition,omitempty"`
	Name           string `toml:"name,omitempty"`
	PackageVersion string `toml:"version,omitempty"`
	RustVersion    string `toml:"rust-version,omitempty"`
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

func (cfg *Metadata) ExtractRustVersionFromCargoFile(fs *afero.Fs) (*string, error) {
	fallbackVersion := "stable"
	f, err := readCargoToml()
	if err != nil {
		return nil, err
	}

	if f.Pkg.RustVersion == "" {
		//nolint:all
		pterm.Warning.Printfln(`no version was retrieved from '%s'.
this is usually because you have not set the 'rust-version'(https://doc.rust-lang.org/cargo/reference/manifest.html#the-rust-version-field) field in your '%s' file.`,
			cfg.ProgrammingLanguageConfigFile,
			cfg.ProgrammingLanguageConfigFile,
		)
		pterm.Warning.Printfln("falling bach to '%s' version", fallbackVersion)
		return &fallbackVersion, nil
	}

	return &f.Pkg.RustVersion, nil
}

func (cfg *Metadata) SetRustVersion(n string, fs *afero.Fs) error {
	if n == "" {
		fmt.Printf(
			"%s version was not set, extracting from '%s' file\n",
			cfg.ProgrammingLanguage,
			cfg.ProgrammingLanguageConfigFile,
		)
		rv, err := cfg.ExtractRustVersionFromCargoFile(fs)
		if err != nil {
			return fmt.Errorf(
				"failed to extract %s version from '%s': %v",
				cfg.ProgrammingLanguage,
				cfg.ProgrammingLanguageConfigFile,
				err,
			)
		}
		pterm.Printfln(
			"%s version extracted from '%s': %s",
			cfg.ProgrammingLanguage,
			cfg.ProgrammingLanguageConfigFile,
			pterm.Yellow(*rv),
		)
		cfg.ProgrammingLanguageVersion = *rv
		return nil
	}
	cfg.ProgrammingLanguageVersion = n
	return nil
}
