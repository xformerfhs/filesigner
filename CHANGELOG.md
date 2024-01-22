# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
and [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

## [0.50.0] - 2024-01-22

### Changed
- Totally new interface for signing.
- Major restructuring of code code modules.

## [0.13.4] - 2023-07-26

### Changed
- Refactored key hashing out of signers and verifiers.

## [0.13.3] - 2023-07-20

### Changed
- Removed unnecessary byte array copy operation.

## [0.13.2] - 2023-06-29

### Changed
- Grammatically correct summary messages.

## [0.13.1] - 2023-06-14

### Changed
- Files can be excluded from being signed

## [0.13.0] - 2023-06-13

### Changed
- Let the user choose the signature type

## [0.12.2] - 2023-06-13

### Changed
- Reintroduced signature type

## [0.12.1] - 2023-06-12

### Changed
- Introduce key id
- Prepare for multiple signature algorithms

## [0.12.0] - 2023-06-12

### Changed
- Back to Ed25519
- Context is now part of the hash

## [0.11.0] - 2023-06-09

### Changed
- Switch to elliptic curve P521

## [0.10.0] - 2023-06-08

### Changed
- Totally reworked sign/verify interface.
- Added maphelper.

## [0.9.8] - 2023-06-06

### Changed
- Glob on Windows should only find files, not directories.

## [0.9.7] - 2023-05-04

### Changed
- Always log public key.

## [0.9.6] - 2023-05-03

### Changed
- Use "constraints" and "maps" packages from "golang/x/exp".

## [0.9.5] - 2023-05-01

### Changed
- Glob works correctly now on Window. Getting the real file name is no longer needed.

## [0.9.4] - 2023-04-26

### Changed
- Get real file names as recorded in the directory. 

## [0.9.3] - 2023-04-21

### Changed
- Made waiting for asynchronous hashers robust. 

## [0.9.2] - 2023-04-20

### Changed
- Time stamp and host name are used for hashing.
- All file paths are sorted alphabetically in output lines.

## [0.9.1] - 2023-04-04

### Added
- New command line format for the "sign" command.

## [0.9.0] - 2023-04-03

### Added
- Initial release, not stable, yet.
