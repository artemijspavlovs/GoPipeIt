package wizard

import (
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"

	"github.com/artemijspavlovs/gopipeit/v2/internal/metadata"
	"github.com/artemijspavlovs/gopipeit/v2/internal/state"
)

func New(metadata *metadata.Metadata, fs afero.Fs) {
	pni := pterm.DefaultInteractiveTextInput
	pni.WithMultiLine(false)
	pn, _ := pni.WithDefaultText("Input a custom project name ( default to the project name defined in your go.mod file )").Show()

	err := metadata.SetProjectName(pn, fs)
	if err != nil {
		pterm.Error.Println("failed to set project name: " + err.Error())
		return
	}

	gbi := pterm.DefaultInteractiveTextInput
	gbi.WithMultiLine(false)
	gb, _ := gbi.WithDefaultText("Input default git branch ( defaults to main )").Show()

	metadata.SetGitBranch(gb)

	gvi := pterm.DefaultInteractiveTextInput
	gvi.WithMultiLine(false)
	gv, _ := gvi.WithDefaultText("Input the Go version to use ( defaults to the Go version defined in your go.mod file )").Show()

	err = metadata.SetGoVersion(gv, fs)
	if err != nil {
		pterm.Error.Println("failed to extract metadata values: ", err.Error())
		return
	}

	cicd := newSelectWizard("Select a CI/CD platform", state.Platforms)
	metadata.SetCICDPlatform(cicd)

	t := newMultiselectWizard("Select tasks to include in the CI/CD pipeline", state.Tasks)
	metadata.SetPipelineTasks(t)

	lt := newMultiselectWizard("Select additional tools that you want to generate the config file for", state.LocalTasks)
	metadata.SetLocalTasks(lt)
}

func newMultiselectWizard(h string, o map[string]state.ConfigurableTool) []string {
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

func newSelectWizard(h string, o map[string]state.ConfigurableTool) string {
	opts := make([]string, 0, len(o))

	for k := range o {
		opts = append(opts, k)
	}

	printer := pterm.DefaultInteractiveSelect.WithOptions(opts)
	printer.DefaultText = h

	selectedOption, _ := printer.Show()

	return selectedOption
}
