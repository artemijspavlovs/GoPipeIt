package wizard

import (
	"errors"
	"fmt"

	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"

	"github.com/artemijspavlovs/gopipeit/internal/metadata"
)

// New function bootstraps together the interactive configurator that the end user sees in their terminals
func New(m *metadata.Metadata, fs *afero.Fs) error {
	pni := pterm.DefaultInteractiveTextInput

	// detect the programming language of the project by the corresponding configuration files
	if ok, _ := afero.Exists(*fs, "go.mod"); ok {
		m.SetProgrammingLanguage("go")
	} else if ok, _ := afero.Exists(*fs, "Cargo.toml"); ok {
		m.SetProgrammingLanguage("rust")
	} else {
		return errors.New(`the tool currently supports only Go(go.mod) and Rust(cargo.toml).
		Neither of the respectful language specific configuration files were find, exiting`)
	}
	pterm.Info.Printfln("language detected: '%s'", m.ProgrammingLanguage)

	pni.WithMultiLine()
	pn, _ := pni.WithDefaultText(
		fmt.Sprintf(
			"Input a custom project name ( defaults to the project name defined in your '%s' file )",
			m.ProgrammingLanguageConfigFile,
		),
	).Show()

	err := m.SetProjectName(pn, fs)
	if err != nil {
		pterm.Error.Println("failed to set project name: " + err.Error())
		return err
	}

	gbi := pterm.DefaultInteractiveTextInput
	gbi.WithMultiLine(false)
	gb, _ := gbi.WithDefaultText("Input default git branch ( defaults to main )").Show()

	m.SetGitBranch(gb)

	gvi := pterm.DefaultInteractiveTextInput
	gvi.WithMultiLine(false)
	gv, _ := gvi.WithDefaultText(
		fmt.Sprintf(
			"Input the %s version to use ( defaults to the %s version defined in your '%s' file )",
			m.ProgrammingLanguage,
			m.ProgrammingLanguage,
			m.ProgrammingLanguageConfigFile,
		),
	).
		Show()

	switch m.ProgrammingLanguage {
	case "go":
		err = m.SetGoVersion(gv, fs)
		if err != nil {
			pterm.Error.Println("failed to extract metadata values: ", err.Error())
			return err
		}
	case "rust":
		err = m.SetRustVersion(gv, fs)
		if err != nil {
			pterm.Error.Println("failed to extract metadata values: ", err.Error())
			return err
		}
	}

	cicd := newSelectWizard("Select a CI/CD platform", metadata.Platforms)
	m.SetCICDPlatform(cicd)

	t := newMultiselectWizard("Select tasks to include in the CI/CD pipeline", metadata.Tasks)
	m.SetPipelineTasks(t)

	lt := newMultiselectWizard(
		"Select additional tools that you want to generate the config file for",
		metadata.LocalTasks,
	)
	m.SetLocalTasks(lt)
	return nil
}

func newMultiselectWizard(h string, o map[string]metadata.ConfigurableTool) []string {
	opts := make([]string, 0, len(o))

	for k := range o {
		opts = append(opts, k)
	}

	printer := pterm.DefaultInteractiveMultiselect.WithOptions(opts)
	printer.DefaultText = h
	printer.Filter = false
	printer.KeyConfirm = keys.Enter
	printer.KeySelect = keys.Space

	selectedOptions, _ := printer.Show()
	return selectedOptions
}

func newSelectWizard(h string, o map[string]metadata.ConfigurableTool) string {
	opts := make([]string, 0, len(o))

	for k := range o {
		opts = append(opts, k)
	}

	printer := pterm.DefaultInteractiveSelect.WithOptions(opts)
	printer.DefaultText = h

	selectedOption, _ := printer.Show()

	return selectedOption
}
