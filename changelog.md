# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0 - 11/11/21 ](https://github.com/ktasper/awsrm/releases/tag/0.0.2)

### Added

- `Makefile` with useful targets
- Moved to the [cobra framework](https://github.com/spf13/cobra)
- Added `s3` as a subcommand via the [cobra framework](https://github.com/spf13/cobra)
- Added `quiet` mode support on the `s3` so you are not prompted to delete buckets.
- A `version` subcommand to show the current version

#### S3

- Will empty and delete S3 buckets.
- Will create a new client and switch to the correct region to delete the buckets you provide.
- Will not delete a S3 bucket if the name is the same as a `vpc` in the same region. (Unless `--skip-vpc-check` is set)

### Other

- Cleaned the imports
- Updated `go.sum`
- Added a release script

## [[0.0.1] - 21-09-21](https://github.com/ktasper/awsrm/releases/tag/0.0.1) 

### Added

- This CHANGELOG file
- Initial `awsrm` that deals with S3 only
- Initial `README.md`