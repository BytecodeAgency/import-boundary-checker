# Valid Go example 3

All packages are not allowed to import any package from within the project.
The domain layer is allowed to import from the data/interactions layer and the
data package is allowed import anything from within the data package and their
subpackages. Since it only imports data/interactions, it does not fail.
