package cmd

import (
	"embed"
	"os"
	"text/template"

	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/artemijspavlovs/gopipeit/internal/metadata"
	"github.com/artemijspavlovs/gopipeit/internal/templates"
	"github.com/artemijspavlovs/gopipeit/internal/wizard"
)

var ApplicationFileSystem = afero.NewOsFs()
var RegenerateAll bool

//go:embed templates/*
var embeddedTemplates embed.FS

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopipeit",
	Short: "Generate CI configuration files with one command",
	Long: `Binary created to provide CI and local development configuration files for Go projects.
Use it to generate optimal configuration files for GitHub Actions, goreleaser, pre-commit and golangi-lint`,
	Run: func(cmd *cobra.Command, args []string) {
		m := metadata.New()

		err := wizard.New(m, ApplicationFileSystem)
		if err != nil {
			return
		}

		r, _ := pterm.DefaultInteractiveConfirm.WithDefaultText("overwrite existing configs? (defaults to No)").Show()
		RegenerateAll = r

		pterm.Info.Println("generating configuration files...")
		tmpl := templates.New()

		pterm.Info.Println("selected CI/CD platform:", pterm.Yellow(m.CICDPlatform))
		switch m.CICDPlatform {
		case "github":
			// move directory creation to CreateDirectoryStructure function?
			err := tmpl.CreateDirectoryStructure(metadata.Platforms["github"].DirectoryStructure, ApplicationFileSystem)
			if err != nil {
				pterm.Fatal.Println("failed to create directory structure for GitHub CI", err)
				return
			}
		}

		if len(m.PipelineTasks) == 0 {
			pterm.Warning.Println("no pipeline tasks were selected, skipping configuration file generation")
		} else {
			pterm.Info.Printfln("generating configuration files related to %s tasks", pterm.Yellow(m.CICDPlatform))
			for t := range m.PipelineTasks {
				err := GenerateConfigFromTemplates(t, metadata.Tasks[t].Configs, m)
				if err != nil {
					pterm.Fatal.Printfln("failed to generate config for %s: %v", t, err)
					return
				}
			}
		}

		if len(m.LocalTasks) == 0 {
			pterm.Warning.Println("no additional tools were selected, skipping configuration file generation")
		} else {
			pterm.Info.Println("setting up additional tools")
			for t := range m.LocalTasks {
				err := GenerateConfigFromTemplates(t, metadata.LocalTasks[t].Configs, m)
				if err != nil {
					pterm.Fatal.Printfln("failed to generate config for %s: %v", t, err)
					return
				}
			}
		}

	},
}

func GenerateConfigFromTemplates(t string, s []metadata.SourceToDest, m *metadata.Metadata) error {
	fs := afero.NewOsFs()
	for _, pair := range s {
		exists, _ := afero.Exists(fs, pair.ConfigDestination)
		if exists && !RegenerateAll {
			pterm.Warning.Printfln(
				"[%s] config %s already exists, it will not be replaced",
				pterm.Yellow(t),
				pair.ConfigDestination,
			)
			continue
		}
		pterm.Info.Printfln("[%s] generating config %s", pterm.Yellow(t), pterm.Yellow(pair.ConfigDestination))

		f, err := fs.Create(pair.ConfigDestination)
		if err != nil {
			return err
		}

		// there are moments when we simply need to create a file as a tool dependency,
		// in this scenario - we simply shouldn't provide a TemplateSource, managed in /internal/state/state.go
		if pair.TemplateSource != "" {
			tmpl := template.Must(template.ParseFS(embeddedTemplates, pair.TemplateSource))

			err = templates.WriteToFile(tmpl, f, m)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gopipeit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().StringVar(
	//	&ProjectName,
	//	"project",
	//	"",
	//	"Project name",
	//)
}
