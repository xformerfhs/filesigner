# filesigner

A program to have an easy way to use digital signatures.
Like certificates, only better and easier ;-).

## Introduction

When one uses artifacts, one needs a method to check whether these artifacts are the ones that the creator created.
Currently there is no way to check if artifacts have been manipulated.
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

How can you now protect against an attacker dropping his forged artifact and generating the signature with the appropriate program?

When artifacts are published, they and the signature file are stored.
The public key used is made known to the recipients of the artifacts by another means.
This allows them to always verify that the signature file is the one issued by the publishing team.

This is illustrated below with an example.

## Calls

### Signing

The signing call looks like this:

```
filesigner sign {contextId} [-type {type}] [-if|-include-file {mask}] [-xf|-exclude-file {mask}] [-id|-include-dir {mask}] [-xd|-exclude-dir {mask}] [-no-subdirs]
```

The parts have the following meaning:

| Part | Meaning |
|----------------------|---------------------------------------------------------------------------------------------------|
| `contextId`          | An arbitrary text used to make the signature depend on a topic, also called a "domain separator". |
| `type`               | Specification of the signature method. Either [`ed25519`](https://en.wikipedia.org/wiki/EdDSA) or `ecdsap521`. If the type is not specified, `ed25519` is used. |
| `include-file`, `if` | Specification of files to include. |
| `exclude-file`, `xf` | Specification of files to exclude. |
| `include-dir`, `id`  | Specification of directories to include. |
| `exclude-dir`, `xd`  | Specification of directories to exclude. |
| `no-subdirs`         | Only process files in the current directory. Do not descend into subdirectories. |

Please note the following information:

* All exclude/include options take one specification.
* Wildcards (`*`, `?`) may be used in include/exclude options.
* An include option excludes all objects that are not included.
* On Linux wildcards need to be put in quotes (`'`) or double quotes (`"`).

The file `signatures.json` is **always** excluded and cannot be signed.

The call creates a file named `signatures.json` which has the following format:

```
{
   "Format": 1,
   "PublicKey": "v5zD877tCK5pk5ZVcj6G4ZVhZXCCJNCX79VkTNtHHg95CZXvHvNT",
   "Timestamp": "2023-06-13 08:28:22 +02:00",
   "Hostname": "MDXN01022044",
   "SignatureType": 1,
   "FileSignatures": {
      "common.go": "N7gGGx2GH2nvzpnzg2HZxjHT4v6zXhVg8PKpXvH5XVhj7jH5jTJcJ65KpV9dHx7JtcJd7T7K9NDpxg248d3pdvJJPxhN7TvDJn8XT32",
      "filesigner": "dNKpPh5hCnKCPCv5KcPJ6Jh6dTGTtKCKXZkngjzK3pg5xk9gxZ4cgC6j2KZTh6ht4x3HK8nvvpChxcDH347543XjJV3vvHx2HV5N45T",
      "filesigner.exe": "3xn8PvVcT27txC5CCn9hV2N5hkPPTng6zJxg9NhGP7D5c4ncndZDXzJhPhkhCc3P8Z69Zn9hNgjH234HvhC3329PC24dXJzd4ctcG4C",
      "main.go": "cgzN2V5TCkkP5C5GNk4D22hh95KXHKVkC3Z6hC4xCKHpGGcK7CkKv6XdH2dzGCphpvHznVPkHdhDTVHkvp658DZDXhp9JgGhZgknn2C",
      "sign_command.go": "5jdKTxVJG2pvTXZdT8kj8tpZg6Zd6xnztPK7zT9D3n9T6CxJz8HNj8DpvHngV9g4hjVgHhc7JTcDkhX854NZpd8Ktzthz5TpxPCzt2T",
      "verify_command.go": "GDP2GGPK5pKxG6ngxd83v55DJTx45Gp9KnN87h4g79P2P57DdNgJCTTvK63VG6cZ7nZjztPhVjj764t42z2zk2Nn25h6xP3VngTJg22"
   },
   "DataSignature": "5PDPVvHnTkd8J2CP4cDpJgXXHJj2kDGCTPxp72t7CXGzpDxH9PT263XXzNpG7p5nhJDdkc3vK7VJV3PKg8K5HznZ2D26CKtzv2cHC2j"
}
```

This is a `json` file with the following fields

| Field             | Meaning                                                                                                                                                                                                                                                                                                          |
|------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Format`         | A number which is the format identifer |
| `PublicKey`      | The public key in a special base32 encoding | 
| `Timestamp`      | The timestamp of the signature | 
| `Hostname`       | The host name of the machine the signature was created | 
| `SignatureType`  | Signature type<br/>          1: [Ed25519](https://en.wikipedia.org/wiki/EdDSA) with [SHA3-512](https://en.wikipedia.org/wiki/SHA-3)-Hash<br/>2: [EcDsa](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm) with curve [P521](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-186.pdf) and [SHA3-512](https://en.wikipedia.org/wiki/SHA-3)-Hash | 
| `FileSignatures` | Signatures of the files as key-value pairs, where the key is the file name and the value is the signature of the file in as special base32 encoding |
| `DataSignature`  | The signature of all the fields in the file in a special base32 encoding |

The signatures are created with the timestamp and the host name.

The return code is:

| Code | Meaning                      |
|------|-------------------------------|
| `0`  | Successful processing     |
| `1`  | Error in the command line   |
| `2`  | Warning while processing |
| `3`  | Error while processing  |

### Verification

The verifcation call looks like this:

```
filesigner verify {contextId} 
```

Here, `contextId` has the same meaning as for the signature.

More parameters are not permitted and will result in an error message.

The program reads the file `signatures.json` and checks whether the files named there exist and whether their signatures match the current content.

The return codes are the same as for signing.

## Programs

| OS      | Program             |
|---------|---------------------|
| Windows | `filesigner.exe` |
| Linux   | `filesigner`     |

The Linux program can be executed on any Linux system.

## Encoding

Binary values are stored in the word-safe [Base32 encoding](https://en.wikipedia.org/wiki/Base32).
What makes this encoding special, is that the alphabet used contains no vowels, no easily confusable characters such as '0' and 'O' or '1' and 'l' and no special characters.
This means that the encoded values can be marked with a double-click, no real words can be created by mistake and no characters can be confused when reading aloud.

## Example

### Signing

Assume that signatures are to be created and checked for the artifacts `filesigner`, `filesigner.exe`, all `Go` files and all files beginning with the word `go` for version `1.7.11` of an application.

The signatures are created with the following call:

```
filesigner sign project1711 -if *.go -if filesign*
```

The program then generates the following output on the console:

```
2023-06-13 12:51:25 +02:00  17  I  filesigner V0.50.0 (go1.21.6, 12 cpus)
2023-06-13 12:51:25 +02:00  37  I  Context id         : project1711
2023-06-13 12:51:25 +02:00  38  I  Public key id      : 85R3-VZPX-JRV8-RN6R-G0L1-SV4U-NW
2023-06-13 12:51:25 +02:00  39  I  Signature timestamp: 2023-06-13 12:51:25 +02:00
2023-06-13 12:51:25 +02:00  40  I  Signature host name: MDXN01022044
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'common.go'
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'filesigner'
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'filesigner.exe'
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'main.go'
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'sign_command.go'
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'verify_command.go'
2023-06-13 12:51:25 +02:00  41  I  Signatures for 6 files successfully created
```

The return code is 0.

### Verifying

To verify the signatures one needs a trusted place where the public key id, the signature timestamp and the signature host name are published.
This may be a signed email, a web site, a database, or whatever is deemed to be a secure trusted place.

Then the verifier runs the filesigner program with the following parameters:

```
filesigner verify project1711
```

The program then generates the following output on the console:

```
2023-06-13 12:51:52 +02:00  17  I  filesigner V0.50.0 (go1.21.6, 12 cpus)
2023-06-13 12:51:53 +02:00  55  I  Context id         : project1711
2023-06-13 12:51:53 +02:00  56  I  Public key id      : 85R3-VZPX-JRV8-RN6R-G0L1-SV4U-NW
2023-06-13 12:51:53 +02:00  57  I  Signature timestamp: 2023-06-13 12:51:25 +02:00
2023-06-13 12:51:53 +02:00  58  I  Signature host name: MDXN01022044
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'common.go'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'filesigner'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'filesigner.exe'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'main.go'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'sign_command.go'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'verify_command.go'
2023-06-13 12:51:53 +02:00  59  I  Verification of 6 files successful
```

The return code is 0.

The verifying person checks, if this ouput if the shown public key id, signature timestamp and signature host are the same as those stored in the trusted place.
If this is not the case, the signature is deemed to be invalid and the files must not be trusted!

As another example, if the file `filesigner` has been manipulated, the following output would appear:

```
2023-06-13 12:51:52 +02:00  17  I  filesigner V0.50.0 (go1.21.6, 12 cpus)
2023-06-13 12:51:53 +02:00  55  I  Context id         : project1711
2023-06-13 12:51:53 +02:00  56  I  Public key id      : 85R3-VZPX-JRV8-RN6R-G0L1-SV4U-NW
2023-06-13 12:51:53 +02:00  57  I  Signature timestamp: 2023-06-13 12:51:25 +02:00
2023-06-13 12:51:53 +02:00  58  I  Signature host name: MDXN01022044
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'common.go'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'filesigner.exe'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'main.go'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'sign_command.go'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'verify_command.go'
2023-06-13 12:54:05 +02:00  22  E  File 'filesigner' has been tampered with
2023-06-13 12:54:05 +02:00  61  E  Verification of 5 files successful and 1 file unsuccessful
```

The return code is 3.

If, for example, the signature file has been manipulated or the wrong context ID has been specified, the following output would appear:

```
2023-06-13 12:51:52 +02:00  17  I  filesigner V0.50.0 (go1.21.6, 12 cpus)
2023-06-13 12:51:53 +02:00  55  I  Context id         : project1711
2023-06-13 12:51:53 +02:00  56  I  Public key id      : 85R3-VZPX-JRV8-RN6R-G0L1-SV4U-NW
2023-06-13 12:51:53 +02:00  57  I  Signature timestamp: 2023-06-13 12:51:25 +02:00
2023-06-13 12:51:53 +02:00  58  I  Signature host name: MDXN01022044
2023-06-13 12:54:56 +02:00  53  E  Signature file has been tampered with or wrong context id
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

## Contakt

Frank Schwab ([Mail](mailto:frank.schwab@live.de "Mail"))

## License

This source code is published under the [Apache License V2](https://www.apache.org/licenses/LICENSE-2.0.txt).
