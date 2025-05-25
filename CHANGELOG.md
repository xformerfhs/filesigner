# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
and [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

## [0.92.0] - 2025-05-25

### Changed
- Added "quiet" option.
- Refactored main program.
- Reordered message numbers.

## [0.91.0] - 2025-05-24

### Changed
- Introduced "verification id".

## [0.83.1] - 2025-05-22

### Changed
- Some internal improvements with no effect on the call interface or output.

## [0.83.0] - 2025-03-01

### Changed
- Correct handling of file errors when verifying files.

## [0.82.2] - 2025-03-01

### Changed
- Signature verification success and error messages all have a severity of "information" now.
- The message numbers for "Not enough arguments", "Context id missing" and "Error parsing command line" have been corrected. They were off by one.

## [0.82.1] - 2025-02-28

### Changed
- Updated dependencies.

## [0.82.0] - 2024-12-23

### Changed
- No longer return an (always nil) error from hash verifier functions.

## [0.81.2] - 2024-10-17

### Changed
- Changed wording of modification/tamper messages.

## [0.81.1] - 2024-08-25

### Changed
- Correct `Set` function `IsProperSubsetOf`.

## [0.81.0] - 2024-08-24

### Changed
- Breaking change: New base32 alphabet that does not produce unwanted character combinations.

## [0.80.1] - 2024-04-23

### Changed
- Made slicehelper.Fill a little faster.

## [0.80.0] - 2024-04-01

### Added
- Initial release.
