package state

type SourceToDest struct {
	//! config destination will depend on the cicd platform
	TemplateSource    string
	ConfigDestination string
}

type ConfigurableTool struct {
	Name               string
	Description        string // TODO: update the wizard to show descriptions of the tools
	DirectoryStructure []string
	Configs            []SourceToDest
}

var Platforms = map[string]ConfigurableTool{
	"github": {
		Name: "GitHub",
		DirectoryStructure: []string{
			".github/workflows",
		},
	},
}

var Tasks = map[string]ConfigurableTool{
	"golangci-lint": {
		Name: "golangci-lint",
		Configs: []SourceToDest{
			{
				TemplateSource:    "templates/github/golangci-lint.yaml.tmpl",
				ConfigDestination: ".github/workflows/golangci-lint.yaml",
			},
			{
				TemplateSource:    "templates/golangci.yaml.tmpl",
				ConfigDestination: "./.golangci.yaml",
			},
		},
	},
	"commitlint": {
		Name:        "commitlint",
		Description: "automatically check whether commits pushed are compatible with conventional commits standard",
		Configs: []SourceToDest{
			{
				TemplateSource:    "templates/github/commitlint.yaml.tmpl",
				ConfigDestination: ".github/workflows/commitlint.yaml",
			},
		},
	},
	"goreleaser": {
		Name:        "goreleaser",
		Description: "automatically release binaries of your application ",
		Configs: []SourceToDest{
			{
				// goreleaser GitHub action
				TemplateSource:    "templates/github/release.yaml.tmpl",
				ConfigDestination: ".github/workflows/release.yaml",
			},
			{
				// goreleaser config
				TemplateSource:    "templates/goreleaser.yaml.tmpl",
				ConfigDestination: "./.goreleaser.yaml",
			},
			{
				// goreleaser prerequisite - CHANGELOG.md file
				ConfigDestination: "./CHANGELOG.md",
			},
		},
	},
	"dependabot": {
		Name: "dependabot",
		Configs: []SourceToDest{
			{
				TemplateSource:    "templates/github/dependabot.yaml.tmpl",
				ConfigDestination: ".github/dependabot.yaml",
			},
		},
	},
}

var LocalTasks = map[string]ConfigurableTool{
	"pre-commit": {
		Name:        "pre-commit",
		Description: "create .pre-commit-config.yaml file with pre",
		Configs: []SourceToDest{
			{
				TemplateSource:    "templates/pre-commit-config.yaml.tmpl",
				ConfigDestination: "./.pre-commit-config.yaml",
			},
		},
	},
	"mkdocs": {
		Name:        "mkdocs",
		Description: "create the configuration file and directory structure for mkdocs framework",
	},
}
