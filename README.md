# gopipeit

`gopipeit` aims to provide a seamless CI/CD setup for a set of tools and language with sensible defaults, as well as CI pipeline configuration files for GitHub Actions (with support for other platforms planned for future releases).

![Alt Text](./docs/gopipeit.gif)

built with [`Cobra`](https://github.com/spf13/cobra) and [`pterm`](https://github.com/pterm/pterm) 🖤

### Generic Tools

- [commitlint](https://github.com/conventional-changelog/commitlint)
- [pre-commit](https://github.com/pre-commit/pre-commit)
- [dependabot](https://github.com/dependabot)
- [mkdocs](https://www.mkdocs.org)

### Go

- [golangci-lint](https://github.com/golangci/golangci-lint)
- [goreleaser](https://github.com/goreleaser/goreleaser)

### Rust

-

### Project Structure

```
cmd/
    templates/ - configuration file templates
    root.go -
internal/
    metadata/ - boostraps the metadata of the project you run `gopipeit` for that is then used to generate the necessary templates
    templates/ - provides
    wizard/ - builds up the cli wizard you see
```
