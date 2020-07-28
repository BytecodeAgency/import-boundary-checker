# Invalid Go example 5

All packages are not allowed to import any package from within the project. The
data package is allowed import anything from within the data package and their
subpackages. The domain/user package imports from data/interactions which is not
allowed for the domain/user package and therefore it fails.
