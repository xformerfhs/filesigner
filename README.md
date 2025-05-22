# filesigner

A program to have an easy way to use digital signatures.
Like certificates, only better and easier ;-).

[![Go Report Card](https://goreportcard.com/badge/github.com/xformerfhs/filesigner)](https://goreportcard.com/report/github.com/xformerfhs/filesigner)
[![License](https://img.shields.io/github/license/xformerfhs/filesigner)](https://github.com/xformerfhs/filesigner/blob/main/LICENSE)

## Introduction

Today's IT world thrives on many, many libraries that are created by volunteers and maintained and developed with a great deal of energy and dedication.
But the users of these libraries are exposed to great danger.
When one uses artifacts, one needs a method to check whether these artifacts are the ones that the creator created.
Currently, there is no way to check if artifacts have been manipulated.
There are several ways to perform such a check:

- Use of hash values
- Use of digital signatures

Hash values are easy to calculate and verify, but they only provide protection against erroneous changes, not against attacks.
An attacker who can change the artifacts can also change the published hash values.

Digital signatures provide protection against such attacks because the attacker would need to have access to the signature's private key to forge it.
They usually work with certificates, but these are difficult to handle:

- It is essential that the private key of a certificate be strongly protected.
If it becomes known to unauthorized persons, they can create signatures for manipulated artifacts themselves.
- It is always necessary to check whether the certificate is valid and whether it has been revoked in the meantime.
This is extremely tedious and error-prone.
- Certificates expire and must be renewed regularly.

This means a considerable organizational effort.

## ATTENTION

The current version still has a `0` in the version number.
So the interface may still change.
Any constructive feedback on this program is welcome.

## Description

The solution presented here provides a digital signature without the hassles associated with certificates.
Artifacts are signed.
To do this, a key pair is generated from a private and public key.
The signatures are created with the private key.
The public key is published so that it can be used to verify the signatures.
After the signing process, the private key is deleted.
It is not stored and therefore cannot be stolen and misused by attackers.
However, it is still possible to verify the signature using the public key.

For a more detailed description of the concept behind this software refer to the [concept.md](doc/en/concept.md) document.

Now, how can you protect against an attacker dropping his forged artifact and generating the signature with the appropriate program?

When artifacts are published, they and the signatures file are stored.
The public key used is made known to the recipients of the artifacts by another means.
This allows them to always verify that the signatures file is the one issued by the publishing team.

This is illustrated below with an example.

## Calls

### Signing

The signing call looks like this:

```
filesigner sign {contextId} [-a|--algorithm {algorithm}] [-i|--include-file {pattern}] [-x|--exclude-file {pattern}] [-I|--include-dir {pattern}] [-X|--exclude-dir {pattern}] [-f|--from-file {file}] [-m|--name {name}] [-r|--recurse] [-s|--stdin] [files...]
```

The parts have the following meaning:

| Part           | Meaning                                                                                                                                                         |
|----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `contextId`    | An arbitrary text used to make the signature depend on a topic, also called a "domain separator".                                                               |
| `algorithm`    | Specification of the signature method. Either [`ed25519`](https://en.wikipedia.org/wiki/EdDSA) or `ecdsap521`. If the type is not specified, `ed25519` is used. |
| `exclude-dir`  | Specification of directories to exclude.                                                                                                                        |
| `exclude-file` | Specification of files to exclude.                                                                                                                              |
| `from-file`    | Read file names to process from the specified file. There is one file name per line.                                                                            |
| `include-dir`  | Specification of directories to include.                                                                                                                        |
| `include-file` | Specification of files to include.                                                                                                                              |
| `name`         | The signatures file name is `{name}-signatures.json`. Default for the name is `filesigner`.                                                                     |
| `recurse`      | Descend also into subdirectories.                                                                                                                               |
| `stdin`        | Read file names to process from the standard input. There is one file name per line.                                                                            |
| `files`        | A blank-separated list of files to sign.                                                                                                                        |

Please note the following information:

* The exclude/include options scan the current directory and the subdirectories if `--recurse` is specified.
* All exclude/include options take one specification.
* Wildcards (`*`, `?`) may be used in include/exclude options.
* An include option excludes all objects that are not included.
* If both, files and includes are specified, they are combined.
* If both, files and excludes are specified, files that match an exclude specification are not processed.
* If wildcards are specified in the files list, they are treated as if they are values in `--include-file` options. 
* On Linux, wildcards need to be put in quotes (`'`) or double quotes (`"`) or escaped by a \\ (like e.g. `--exclude-dir .\*` to exclude all directories starting with `.`).

> [!IMPORTANT]
> The signatures file is **always** excluded and cannot be signed.

The call creates a signatures file[^1] which has the following format:

[^1]: The content of the file is just an illustration.
  It is not possible to verify the signatures of files in this repository with it.

```
{
   "format": 1,
   "contextId": "project1711",
   "publicKey": "bfdjDDJ44djrcjhRfFRtdz4HFJjjdZTzBSm37FrZjDgMzFjdv7zT",
   "timestamp": "2024-03-05 15:48:51 +01:00",
   "hostname": "BuildHost",
   "signatureType": 1,
   "fileSignatures": {
      "common.go": "3sLc4CVCsMmfgmhtf4LBssGt9rtrZHmJhRrVB3QQ7M7LRdCvjGHh3rdHDH77mQgFC3Z9f9jmcDjdtRDFGS9QgC37r3QrHZSfzvcZ743",
      "filesigner": "FsVSjJSbQTVLgfSJLvRS33G9bFHdMFSRFGb9FL94C9v37D7zrZSCBmRrFhDfQcHGTFhbhjFFZRhQVMJ4sGB3FVFSdfDgtfRgBzDLJ9G",
      "filesigner.exe": "HMQm44ShmbLcfQSv94vcGsMHZLJ4VVZMsfgcbHDVDcbtQg4RTjBCBsm9b94rgDVLCgQdD4GHbBLzFM7RTGhQB94MQ9HQvMgRcQT7Q4h",
      "filesigner_sbom.json": "bDvRzQG9dLDTQGVBHfrzJfJBTrCBDhrzbMVvsc7FbbzcFhM7FGtLzLftBCL9fVzRFrgMDMcCsmCTtTQc7j4fBmSGR7rfBSrs9tbbB4G",
      "main.go": "QjgzMhsRSSsJMjBDhfTVm3BmBdSZzMZCTvbG7TssCrDVHG9mVrhGHjMVdvdthrrLrdr4jCJbZDfGstJsrdCSJR4gtSRMg9fSzbCrM9h",
      "sign_command.go": "HcDcF7LmB4mSvvVJfTvbSSgCtcc7t3vCcjCjZgQc3jfGJ3MfdSS9FChQV37LjdBVhMChLsrdv9vQSJbmgSfD9HszhVJhSDdZL4TdM33",
      "verify_command.go": "jQMcCcjRTVQmM7StcZbfVZmfJbstLv9FSFSDZrdrsVVBbHdJQ74BbcH4hsz4gL7V9cvTzgFSZDhVtBMhzbmDFRdR4Sg94ZSVDm7Qc7h"
   },
   "dataSignature": "cRSSgg3Cbhd3vGLgBf4MtFDCHb3SZrtzCJDrZmhHMzrHVVd3Vc94zhJLmD9g7CRdSRjhjQQhhszCmd3rH79LGT3t97tTg4Gt7bCHc7T"
}
```

This is a `json` file whose structure and the meaning of the fields are described in [file_format.md](doc/en/file_format.md).
A detailed description of the various calculations and data formats can be found in [technical_specification.md](doc/en/technical_specification.md).

The possible return codes are the following:

| Code | Meaning                   |
|------|---------------------------|
| `0`  | Successful processing     |
| `1`  | Error in the command line |
| `2`  | Warning while processing  |
| `3`  | Error while processing    |

### Verification

The verification call looks like this:

```
filesigner verify [-m|--name {name}]
```

The parts have the following meaning:

| Part           | Meaning                                                                                                                                                         |
|----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`         | The signatures file name is `{name}-signatures.json`. Default for the name is `filesigner`.                                                                     |

> [!IMPORTANT]
> More parameters are not permitted and will result in an error message.

The program reads the signatures file and checks whether the files named there exist and whether their signatures match the current content.

The return codes are the same as for signing.

## Programs

| OS      | Program          |
|---------|------------------|
| Windows | `filesigner.exe` |
| Linux   | `filesigner`     |

The Linux program can be executed on any Linux system.

## Encoding

Binary values are stored in a special [Base32 encoding](https://en.wikipedia.org/wiki/Base32).
What makes this encoding special, is that the alphabet used contains no vowels, no easily confusable characters such as '0' and 'O' or '1' and 'l' and no special characters.
This means that the encoded values can be marked with a double click, no real words can be created by mistake and no characters can be confused when reading aloud.

## Example

### Signing

Assume that signatures are to be created and checked for the artifacts `filesigner`, `filesigner.exe`, all `Go` files and all files beginning with the word `go` for version `1.7.11` of an application.

The signatures are created with the following call:

```
filesigner sign project1711 -if *.go -if filesign*
```

The program then generates the following output on the console:

```
2024-03-05 15:48:51 +01:00  15  I  filesigner V0.83.1 (go1.24.3, 8 cpus)
2024-03-05 15:48:51 +01:00  24  I  Context id         : project1711
2024-03-05 15:48:51 +01:00  25  I  Public key id      : DLQB-J6MT-YMF1-PPRF-KQ6P-V9LG-QR
2024-03-05 15:48:51 +01:00  26  I  Signature timestamp: 2024-03-05 15:48:51 +01:00
2024-03-05 15:48:51 +01:00  27  I  Signature host name: Jetzt
2024-03-05 15:48:51 +01:00  21  I  Signing succeeded for file 'common.go'
2024-03-05 15:48:51 +01:00  21  I  Signing succeeded for file 'filesigner'
2024-03-05 15:48:51 +01:00  21  I  Signing succeeded for file 'filesigner.exe'
2024-03-05 15:48:51 +01:00  21  I  Signing succeeded for file 'filesigner_sbom.json'
2024-03-05 15:48:51 +01:00  21  I  Signing succeeded for file 'main.go'
2024-03-05 15:48:51 +01:00  21  I  Signing succeeded for file 'sign_command.go'
2024-03-05 15:48:51 +01:00  21  I  Signing succeeded for file 'verify_command.go'
2024-03-05 15:48:51 +01:00  37  I  Signatures for 7 files successfully created and written to 'filesigner-signatures.json'
```

The return code is 0.

### Verifying

To verify the signatures one needs a trusted place where the public key id, the signature timestamp and the signature host name are published.
This may be a signed email, a website, a database, or whatever is deemed to be a secure trusted place.

Then the verifier runs the filesigner program with the following parameters:

```
filesigner verify
```

The program then generates the following output on the console:

```
2024-03-05 15:49:13 +01:00  15  I  filesigner V0.83.1 (go1.24.3, 8 cpus)
2024-03-05 15:49:13 +01:00  51  I  Reading signatures file 'filesigner-signatures.json'
2024-03-05 15:49:13 +01:00  24  I  Context id         : project1711
2024-03-05 15:49:13 +01:00  25  I  Public key id      : DLQB-J6MT-YMF1-PPRF-KQ6P-V9LG-QR
2024-03-05 15:49:13 +01:00  26  I  Signature timestamp: 2024-03-05 15:48:51 +01:00
2024-03-05 15:49:13 +01:00  27  I  Signature host name: Jetzt
2024-03-05 15:49:13 +01:00  21  I  Verification succeeded for file 'common.go'
2024-03-05 15:49:13 +01:00  21  I  Verification succeeded for file 'filesigner'
2024-03-05 15:49:13 +01:00  21  I  Verification succeeded for file 'filesigner.exe'
2024-03-05 15:49:13 +01:00  21  I  Verification succeeded for file 'filesigner_sbom.json'
2024-03-05 15:49:13 +01:00  21  I  Verification succeeded for file 'main.go'
2024-03-05 15:49:13 +01:00  21  I  Verification succeeded for file 'sign_command.go'
2024-03-05 15:49:13 +01:00  21  I  Verification succeeded for file 'verify_command.go'
2024-03-05 15:49:13 +01:00  56  I  Verification of 7 files successful
```

The return code is 0.

The verifying person checks, if the shown public key id, signature timestamp and signature host are the same as those stored in the trusted place.
If this is not the case, the signature is deemed to be invalid and the files must not be trusted!

As another example, if the file `filesigner` has been manipulated, the following output would appear:

```
2024-03-05 15:49:38 +01:00  15  I  filesigner V0.83.1 (go1.24.3, 8 cpus)
2024-03-05 15:49:38 +01:00  51  I  Reading signatures file 'filesigner-signatures.json'
2024-03-05 15:49:38 +01:00  24  I  Context id         : project1711
2024-03-05 15:49:38 +01:00  25  I  Public key id      : DLQB-J6MT-YMF1-PPRF-KQ6P-V9LG-QR
2024-03-05 15:49:38 +01:00  26  I  Signature timestamp: 2024-03-05 15:48:51 +01:00
2024-03-05 15:49:38 +01:00  27  I  Signature host name: Jetzt
2024-03-05 15:49:38 +01:00  21  I  Verification succeeded for file 'common.go'
2024-03-05 15:49:38 +01:00  21  I  Verification succeeded for file 'filesigner.exe'
2024-03-05 15:49:38 +01:00  21  I  Verification succeeded for file 'filesigner_sbom.json'
2024-03-05 15:49:38 +01:00  21  I  Verification succeeded for file 'main.go'
2024-03-05 15:49:38 +01:00  21  I  Verification succeeded for file 'sign_command.go'
2024-03-05 15:49:38 +01:00  21  I  Verification succeeded for file 'verify_command.go'
2024-03-05 15:49:38 +01:00  22  E  File 'filesigner' has been modified
2024-03-05 15:49:38 +01:00  58  I  Verification of 6 files successful and 1 file unsuccessful
```

The return code is 3.

If, for example, the signatures file has been manipulated the following output would appear:

```
2024-03-05 15:50:04 +01:00  15  I  filesigner V0.83.1 (go1.24.3, 8 cpus)
2024-03-05 15:50:04 +01:00  51  I  Reading signatures file 'filesigner-signatures.json'
2024-03-05 15:50:04 +01:00  54  E  Signatures file has been modified
```

The return code is 3.

## Program build

You must have Go installed to create the program.
This creates a directory that is specified in the `GOPATH` environment variable.
Under Windows, this is the home directory, e.g. `D:\Users\username\go`.
Under Linux this is `${HOME}/go`.
In that directory there is a subdirectory `src`.

To create the program, the source code must be stored under `%GOPATH%\src\filesigner` or `${HOME}/go/src/filesigner`.
Then one has to start the batch file `gb.bat` or the shell script `gb`, which builds the executables.
These scripts expect the UPX program to be in a specific location.
This location can be adapted to the local path.
If UPX is not available, no compression is performed.

As a result, the files `filesigner` for Linux and `filesigner.exe` for Windows are created.

## Contact

Frank Schwab ([Mail](mailto:github.sfdhi@slmails.com "Mail"))

## License

This source code is published under the [Apache License V2](https://www.apache.org/licenses/LICENSE-2.0.txt).
