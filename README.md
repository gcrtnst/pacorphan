# pacorphan

pacorphan is a pacman utility that helps you find packages no longer needed by any explicitly installed package.

## How pacorphan differs from `pacman -Qdt`

`pacman -Qdt` simply reports packages that are "not required by any other package". That works well for simple cases, but it only surfaces the terminal packages of the dependency graph. Packages with dependents are not detected, even if they are not required by any explicitly installed packages. Packages involved in circular dependencies are also not detected properly.

In contrast, pacorphan walks the dependency graph from explicitly installed packages and reports all packages that are not directly or indirectly depended on by them. This makes it possible to:

- find orphan packages that sit above other dependencies in the chain
- find orphan packages that are part of circular dependency relationships
- list all orphan packages at once

### Dependency example

A dependency chain can look like this:

```mermaid
flowchart LR
    A[package A]
    A --> B[package B]
    B --> C[package C]
```

If none of the packages are explicitly installed, `pacman -Qdt` will report only `A`, but `pacorphan` will report all of `A`, `B`, and `C`.

Circular dependencies are another case where pacorphan is more useful:

```mermaid
flowchart LR
    X[package X]
    Y[package Y]
    X --> Y
    Y --> X
```

If neither package is explicitly installed, pacorphan can detect the cycle as orphaned. `pacman -Qdt` does not handle this scenario well.

## Installation

Install it with Go:

```sh
go install github.com/gcrtnst/pacorphan@latest
```

Make sure `go` and `base-devel` are installed before running the command.

Ensure that `$(go env GOPATH)/bin` is added to your PATH environment variable.

## Usage

List orphan packages:

```sh
pacorphan
```

Show package names only:

```sh
pacorphan -q
```

By default, pacorphan includes optional dependencies in the analysis. Use `-t` to ignore them:

```sh
pacorphan -t
```

A complete list of command-line flags is available with:

```sh
pacorphan --help
```

### Removing orphans

> [!WARNING]
> Review the list carefully before removing packages.

> [!NOTE]
> You don't need to add `-s` (`--recursive`) to `pacman -R`; pacorphan already handles that.

Remove all orphan packages:

```sh
pacorphan -q | sudo pacman -R -
```

To also remove configuration files without creating backups, add `-n` to pacman:

```sh
pacorphan -q | sudo pacman -Rn -
```

To ignore optional dependencies and remove configuration files without backups:

```sh
pacorphan -qt | sudo pacman -Rn -
```

## License

This project is licensed under the terms of the [LICENSE](LICENSE) file.
