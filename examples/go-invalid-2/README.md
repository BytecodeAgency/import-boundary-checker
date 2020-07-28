# Invalid Go example 2

The domain layer is not allowed to import from the data layer but from the
data/interactors layer. Since it only import from data/database, it fails.
