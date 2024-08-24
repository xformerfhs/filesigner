# File format

## Description

This document describes the structure of the signature file.

The data is stored in [JSON](https://en.wikipedia.org/wiki/JSON) format.
All texts are encoded in [UTF-8](https://en.wikipedia.org/wiki/UTF-8).

The following fields are present in the file:

| Field            | Meaning                                                                                                                           |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------|
| `contextId`      | The context id of the signature.                                                                                                  |
| `dataSignature`  | The signature over the individual parts of this file.                                                                             |
| `fileSignatures` | The list of signatures of the individual files as key-value pairs, where the key is the file path and the value is the signature. |
| `format`         | The identifier for the format of this file.                                                                                       |
| `hostname`       | The name of the machine where the signatures were created.                                                                        | 
| `publicKey`      | The public key.                                                                                                                   |
| `signatureType`  | The signature type.                                                                                                               | 
| `timestamp`      | The timestamp of the signature.                                                                                                   | 

### Format identifier

The format identifier specifies the format of the file.
Currently only one value is defined: `1`.
This value means that the file has the structure described here.

### Signature type

The signature type can have two values:

| Signature type | Meaning                                                                                                                                                                                                                                             |
|:--------------:|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|      `1`       | The signatures are created with the algorithm [Ed25519](https://en.wikipedia.org/wiki/EdDSA#Ed25519).                                                                                                                                               |
|      `2`       | The signatures are created using the ECDSAP521 algorithm, i.e. [ECDSA](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm) with the curve [P-521](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-186.pdf). |

Both methods use elliptical curves.
Further information can be found in the file [technical specification.md](technical_specification.md).

### Timestamp

The timestamp is in [ISO 3339](https://datatracker.ietf.org/doc/html/rfc3339) format: `YYYY-MM-DD hh:mm:ss +hh:mm`.

The time stamp begins with the date according to the [Gregorian calendar](https://en.wikipedia.org/wiki/Gregorian_calendar) in the form year (four digits), month (two digits), day (two digits).
The individual parts are separated from each other by a minus sign (`-`).
All parts have leading zeros.
For example, the month of March is written as `03`.

This is followed by a space as the separator.

This is followed by the time in the form hour (two digits, 24 hours), minute (two digits), second (two digits).
The individual parts are separated from each other by a colon (`:`).
All parts have leading zeros.

This is followed by a space as the separator.

This is followed by the time zone in the form sign (`+`, `-`), hours (two digits), a colon (`:`), minutes (two digits).
The value indicates the time difference between the time zone to which the specified timestamp refers and [UTC](https://en.wikipedia.org/wiki/Coordinated_Universal_Time).
Here are a few examples of the value of the time zone specification:

| Place             | Season               | Time zone |
|-------------------|----------------------|-----------|
| London            | Standard time        | `+00:00`  |
| London            | Daylight saving time | `+01:00`  |
| Frankfurt am Main | Standard time        | `+01:00`  |
| Frankfurt am Main | Daylight saving time | `+02:00`  |
| Los Angeles       | Standard time        | `-08:00`  |
| Los Angeles       | Daylight saving time | `-07:00`  |
| Mumbai            | Standard time        | `+05:30`  |

### Encoding of binary data

The public key, the file signatures and the signature of the signature file are binary values that are encoded similar to the [word-safe Base32 method](https://en.wikipedia.org/wiki/Base32#Word-safe_alphabet).

In the used Base32 coding, a character represents 5 bits from the binary value.
The alphabet used is: `3479BCDFGHJLMRQSTVZbcdfghjmrstvz`.
The `3` stands for `0`, the `z` for `31` and the other characters for the values in between.

It does not contain any vowels, so that words that have a meaning are not inadvertently created.
Likewise, there is neither `6` nor `X`or `x` whose repeated stringing together is undesirable.
It also does not contain any special characters, so that any value coded in this way can be completely marked in an editor with a double click.

### File paths

The file paths in the `fileSignatures` field are relative file paths.
They always refer to the current directory.
The slash (`/`) is always used as the path separator, regardless of the operating system used.
The file paths are encoded in UTF-8.

## Rules

The JSON file must be checked for formal errors when it is read in.
The following rules apply:

- All fields **must** be present.
- There **must** be no additional fields.

If at least one field is missing or at least one additional field is present, processing is aborted.

After reading the file, the signature `dataSignature` of the entire file is checked first.
If the overall signature does not match the remaining data in the file, processing is aborted.

## Example

In the following example the above descriptions are explained with specific values .

Let's assume that on 25.02.2024 at 13:37:22 on the computer "BuildHost", which is located in Mumbai, a signatures file was created for the following files:
`common.go`, `filesigner`, `filesigner.exe`, `maphelper/map_helper.go` and `set/set.go`.

The signature file would then look something like this:

```
{
   "format": 1,
   "publicKey": "HxVJVrrjQcgfvhPxJ45chrQrRCFWmgJ5JH8JGMv6xxj23xjH8P52",
   "timestamp": "2024-02-25 13:37:22 +05:30",
   "hostname": "BuildHost",
   "signatureType": 1,
   "fileSignatures": {
      "common.go": "mpmQrxWxqgwPmw54gm6hPMMv8pW4MFRjmJrfqH99H4q39vQgxVJ3vcj2qpcWPqVr337VRPJJ4X3CP8WR39VwgJQjJfW23RPWmhPQW2R",
      "filesigner": "WWQmj6822QrqcwP2j75pw2xf3Hvjpv97H2cXQqVx97WmGP77c25H28PQcjCCQv9MP8xp2cq46R9prQpqGVrGRCVV7gQJw6Q3r3cpp3R",
      "filesigner.exe": "3W9rJ3WGx6mVFFXjXQmMvJ4jw9XhQ7j3VGmhcRj2pF42JgGvP2frhQXg6V5QQvwPqxmcx7Wg86C7v3v4H9qWqPwCMcvpgxxJv3M7442",
      "maphelper/map_helper.go": "HXqJRpVH69Xpxgp9f9FpWMXQxF5rJPPGwFC2Hcp9hG979f4FR8Fgv6fG4mPqQxpW9x6R9RRvWXgGF5rpHW83gCGQr2hPfHMJ6G34J52",
      "set/set.go": "fvw537rpJq3QHHgqxFW8J85VCvjcgFrwRh8gH5wPh266f65V8XmjqGJG9pJ98xm9vMGfFp2jpr6qw4cv46G9vPX5J8F4cR9vG3M8R3R",
   },
   "dataSignature": "PCvCgQ9PFjf6hPCh4RvPHvCp47VghrmX96fwC2r43VhxJCXHXR4v2QX8wrwFjQRm4FqG56cM9wf4pf4hhMrR2JpM88Pf8pFxvmRhj2j"
}
```

The values are easy to understand with the explanations above.

The decoded value for the public key is given here as an example.
The Base32-encoded value `HxVJVrrrjQcgfvhPxJ45chrQrRCFWmgJ5JH8JGMv6xxj23xjH8P52` corresponds to the following byte values: `5f e2 c8 f3 98 7d 2d 5e dd df 60 87 4b f1 fc 82 13 2c d9 83 62 cc c5 37 a4 ff f0 00 ff 0b 33 86`.
