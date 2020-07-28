# Valid Go example 2

The domain layer is not allowed to import from the data layer but from the
data/interactions layer. Since it only import from data/interactions, it does
not fail.
