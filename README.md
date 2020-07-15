# Import Boundary Checker for Golang

[![pipeline status](https://git.bytecode.nl/foss/import-boundry-checker/badges/master/pipeline.svg)](https://git.bytecode.nl/foss/import-boundry-checker/-/commits/master)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=BytecodeAgency_import-boundary-checker&metric=alert_status)](https://sonarcloud.io/dashboard?id=BytecodeAgency_import-boundary-checker)
[![Maintainability](https://api.codeclimate.com/v1/badges/4870797be10646d8ddd0/maintainability)](https://codeclimate.com/github/BytecodeAgency/import-boundary-checker/maintainability)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/BytecodeAgency/import-boundary-checker)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/BytecodeAgency/import-boundary-checker)
![GitHub](https://img.shields.io/github/license/BytecodeAgency/import-boundary-checker)

![import boundary checker](https://github.com/BytecodeAgency/import-boundary-checker/raw/master/examples/examples-go.gif)

## About

Import Boundary Checker is a tool to automatically check if import boundaries are violated or not. Examples of where this tool is useful:

* **Hexagonal/clean architecture** where domain logic cannot import infrastructure code
* **Mono-repositories with multiple microservices** where you don't allow services to import code from other microservices
* **Layered architecture** where layers cannot import certain other layers

Why is this tool useful, and why consider using it over alternative tools?

* Import boundaries are checked **automatically** (meaning you don't have to spend time in code review on manually checking)
* It is extremely **fast** (sub-second speeds with medium sized projects)
* Configuration is **easy** (within a few minutes you can define forbidden imports for your projects and have everything set up)
* The tool is **independent** (you don't need to change any production or testing code for the tool to work)
* Usage of the tool requires **zero dependencies** (so no outside dependencies are needed for running the tool, only the source code of your project)

## Installation

You will need Go 1.14+.

```sh
go install github.com/BytecodeAgency/import-boundary-checker
```

_Note: due to being in the alpha phase, the tool has only been tested on MacOS and Linux. Windows support is not guaranteed but will be added before the 1.0.0 release._

## Usage

After setting your configuration and, you can invoke the tool by simply calling:

```sh
import-boundary-checker
```

Note that only production code (not the test modules) is tested.

You can use the following CLI options (they are all optional):

* `-config [filepath]`: set the configuration path (defaults to `.importrules` in the current directory)
* `-verbose`: set verbose output of main (does not enable debugging mode)

## Configuration DSL

The tool is configured using a domain specific language. Follow the steps below to create your configuration:

0. Create a new configuration file `.importrules` in the root directory of the project you want to check

1. Set the correct language for your project (currently only Go is supported, Typescript/Javascript will be added next)

```
LANG "Go";
```

2. Set the `IMPORTBASE` variable (this is the same as the `module` value in `go.mod`)

```
IMPORTBASE "github.com/BytecodeAgency/example";
```

3. Define the import boundaries, using
    `IMPORTRULE "[IMPORTBASE]{file you are defining forbidden imports for} CANNOTIMPORT "[IMPORTBASE]/{some module in project}" "[IMPORTBASE]/{another module in project}";`.
    Leaving out the `[IMPORTBASE]` allows you to define forbidden imports from the standard library or outside dependencies. Whitespace is ignored.

```
IMPORTRULE "[IMPORTBASE]/typings/entities"
CANNOTIMPORT "[IMPORTBASE]" "fmt" "github.com/go-playground/validator/v10";

IMPORTRULE "[IMPORTBASE]/domain"
CANNOTIMPORT
    "[IMPORTBASE]/infrastructure"
    "[IMPORTBASE]/data"
```

**This will give the following configuration file**:

```
LANG "Go";
IMPORTBASE "github.com/BytecodeAgency/example"

IMPORTRULE "[IMPORTBASE]/typings/entities"
CANNOTIMPORT "[IMPORTBASE]" "fmt" "github.com/go-playground/validator/v10";

IMPORTRULE "[IMPORTBASE]/domain"
CANNOTIMPORT
    "[IMPORTBASE]/infrastructure"
    "[IMPORTBASE]/data"
```

You can read the full DSL specification in the [`docs/dsl.md`](docs/dsl.md) file.

## Full examples

In the [`examples`](/examples) directory of this repository, you can find some examples of the tool configured. In each directory, you will find a `.importrules` file with the import boundaries defined.

An example for running this in CI through Docker can be found in the [`Dockerfile`](Dockerfile) in the root of this repository.

## License and contribution

The project is licensed under the Lesser GNU Public Licence version 3 (LGPLv3). Contributions (issues and pull requests) are welcome.

Development documentation can be found in [`docs/dev.md`](docs/dev.md).
