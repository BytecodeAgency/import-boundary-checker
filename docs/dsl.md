# DSL definition

* Only one language at a time is supported
* Statements should be closed with a semicolon (no implicit statement endings like in Go or JS)
* Language definition uses the following syntax

```
LANG "[Typescript/Go]";
```

* Defining import rules uses the following syntax

```
IMPORTRULE "path/to/module" CANNOTIMPORT "/path/to/other/module";
```

* Multiple forbidden imports can be specified by using whitespace as separator

```
IMPORTRULE "path/to/module" CANNOTIMPORT "/path/to/other/module1" "/path/to/other/module2";
```
* Line endings can be used to make the file more readable as you please (whitespace is ignored)

```
IMPORTRULE "path/to/module"
CANNOTIMPORT
    "/path/to/other/module1"
    "/path/to/other/module2"
    "/path/to/other/module3";
```

* When defining an import path, all sub modules/directories are automatically included, no `/*` wildcard has to be added (the second import statement is redundant)

```
IMPORTRULE "path/to/module"
CANNOTIMPORT
    "/path/to/other/module"
    "/path/to/other/module/sub";
```

## Wishlist for configuration

* Support `DIRECTORY "src";` to define in which directories to run the application (where to start the file/dir walker)
* Support `EXTENSIONS "tsx" "ts"` to define which extensions to include when checking
* Support comments
* Support multiple entries for `IMPORTRULE`
* Support `CANNOTIMPORT "*";` or `CANNOTIMPORT;` to never allow any imports
* Support "CANONLYIMPORT" (useful for domain/entity layers)
* Support exclusions for rules (usecase: prohibit `/some/module` but allow `/some/module/mocks`)
* Support Regex in definitions
* Support configuration of stdlib import rules

## Example config

```
LANG "Go";

IMPORTRULE "git.bytecode.nl/single-projects/youngpwr/platform-backend/typings/entities"
CANNOTIMPORT "git.bytecode.nl/single-projects/youngpwr/platform-backend";

IMPORTRULE "git.bytecode.nl/single-projects/youngpwr/platform-backend/domain"
CANNOTIMPORT
    "git.bytecode.nl/single-projects/youngpwr/platform-backend/infrastructure"
    "git.bytecode.nl/single-projects/youngpwr/platform-backend/data";
```
