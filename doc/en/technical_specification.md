# Technical specification

## General

This file describes the technical details of the calculation and verification of signatures.
The first part explains the cryptographic basics.
The second part specifies the individual calculation steps.

## Cryptographic basics

This system works with digital signatures.
The following parts are required to create a digital signature:

1. A cryptographically secure [hash method](https://en.wikipedia.org/wiki/Cryptographic_hash_function).
2. An [asymmetric encryption method](https://en.wikipedia.org/wiki/Public-key_cryptography) such as [RSA](https://en.wikipedia.org/wiki/RSA_(cryptosystem)) or [elliptic curves](https://en.wikipedia.org/wiki/Elliptic-curve_cryptography).

### Signature

To create a signature a cryptographically secure hash method is used to calculate a cryptographically secure [checksum](https://en.wikipedia.org/wiki/Checksum) from the data to be signed.
This cryptographically secure hash is then encrypted with the **private** - i.e. secret - key of an asymmetric encryption method.
This encrypted hash value is the digital signature.

An asymmetric encryption method uses a key pair consisting of a private and a public key.
The following rules apply with regard to encryption and decryption:

- What is __encrypted__ with the **public** key can only be __decrypted__ with the **private** key.
- What is __encrypted__ with the **private** key can only be __decrypted__ with the **public** key.

I.e. a message encrypted with **one** key can only decrypted with the **other** key of the key pair.

As the name suggests, the public key is known and is not kept secret.
The private key must be kept secret.

### Verification

For verification the verifier calculates the hash value of the received data using the same hash procedure as the creator of the data.
Then the signature is decrypted with the **public** - i.e. known - key of the creator.
The decrypted hash value from the signature is compared to the hash value that has been calculated locally.
If the two hash values match, the digital signature is valid, otherwise it is not.

### Meaning

If the verification is successful, it is ensured that the creator of the signature has used the private key that matches the public key.

The security of this procedure rests on two pillars:

1. The hash procedure ensures that an attacker cannot change the data in such a way that the same hash value is created with the changed data.
2. The encryption method ensures that the value of the private key cannot be deduced from knowledge of the public key.

### Selected methods

The SHA-3-512 method is used to calculate the hash value, i.e. [SHA-3](https://en.wikipedia.org/wiki/SHA-3) with a hash length of 512 bits (64 bytes).
This method was standardized by [NIST](https://www.nist.gov/) and is currently the most secure hash method with a very long and therefore still secure hash value length in the long term.

The signature methods used are [Ed25519](https://en.wikipedia.org/wiki/EdDSA#Ed25519) and [ECDSA](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm) with the curve [P-521](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-186.pdf).

Both methods use elliptic curves as an asymmetric encryption method.
Elliptic curves are currently and in the foreseeable future secure against attacks by classical computers.
Theoretically, they can be attacked by quantum computers.
However, elliptic curves are more difficult to attack by quantum computers than RSA, as the attack has a higher complexity.
If a concrete threat to elliptic curves becomes known, the system can be changed at any time so that a quantum computer-safe signature system such as [CRYSTALS-Dilithium](https://pq-crystals.org/dilithium/), [FALCON](https://falcon-sign.info/) or [SPHINCS+](https://sphincs.org/) is used.
At present, however, there are hardly any libraries available for the use of these systems.

At the same time, elliptic curves are more efficient in terms of resource consumption due to their significantly shorter key lengths compared to RSA.
In summary, there are currently no known effective attacks on elliptic curves.

The Ed25519 method is based on the elliptic curve [Curve25519](https://www.rfc-editor.org/rfc/rfc8032#page-9), which is one of the most secure known elliptic curve and can be computed efficiently.

The ECDSA method with the P-521 curve, on the other hand, uses the P-521 curve standardized by NIST with the ECDSA standard also created by NIST.
The method is significantly less efficient than Ed25519, but is available in many programming languages.

## Calculations

The following section describes how the individual calculations are performed.

Whenever the content of bytes is specified, these are given in hexadecimal notation.

Two values are used several times so that they are described here in advance:

### Counter of variable length

Counters or lengths are used for the calculation.
Their values are always used with as many bytes as necessary and no more.
The bytes are specified in the so-called [big-endian](https://de.wikipedia.org/wiki/Byte-Reihenfolge#Big-Endian-Format) format.

To illustrate this, here are a few examples of values with their encoding in variable length:

|    value | encoding (hexadecimal) |
|---------:|------------------------|
|      `0` | `00`                   |
|      `1` | `01`                   |
|    `255` | `ff`                   |
|    `300` | `01 2c`                |
|  `65432` | `ff 98`                |
| `100000` | `01 86 a0`             |

### Context key

The signatures need a "context id".
This context id is put into the hash calculation of every file.
However, the context id is user-supplied data.
An attacker could specify a malicious context id that could enable attacks[^1] on the signature algorithms.
In order to make such attacks impossible a context key is derived from the context id.
This context key is used in the hash calculations, not the context id.

[^1]: Currently no such attacks are known.

The context key is calculated from the context is as follows:

1. The characters of the context id are encoded in [`UTF-8`](https://en.wikipedia.org/wiki/UTF-8).
2. The length is appended in variable length encoding. This is referred to as the "extended context id".
3. Then the byte sequence of the extended context id is reversed and the SHA-3-256 value of this reversed extended context id is calculated.
4. From this hash with a length of 32 bytes a 64 byte key with the following parts is created:
    - Constant byte sequence `6f 00 11 21 3d 31 c2 3b c3 69 ab 0b 6d 8e 42 35`.
    - Hash value just calculated.
    - Constant byte sequence `30 2d 15 d7 37 d5 b1 df 45 ee 30 bc e0 0b 89 cc`.
5. This 64 byte key is used as the key for an [SHA-3-512-HMAC](https://en.wikipedia.org/wiki/HMAC)-SHA-3-512.
6. The SHA-3-512-HMAC value of the context id bytes is calculated with this key.
7. This generates a 64 byte HMAC value which is used as a padding to create the context key as follows:
    - The first 32 bytes of the HMAC value.
    - The extended context id.
    - The last 32 bytes of the HMAC value.

This context key is used in the calculation of all hash values.
    - First, the first half of the key is fed into the hash value.
    - Then all the data is fed into the hash value.
    - Finally, the second half of the key is fed into the hash value.
    - If the key length is odd, the first half of the key is one byte shorter than the second half.

An example shows this calculation rule using concrete values:

- The context id is `Überführung`.
- Step 1: The context bytes in UTF-8 encoding are `c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67` and have the length 13 (`0d`).
- Step 2: The extended context ID has the value `c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67 0d`.
- Step 3: The SHA-3-256 value of the reversed extended context id (`0d 67 6e 75 72 68 bc c3 66 72 65 62 9c c3`) is calculated which yields `86 3a fd 35 1e 70 d5 07 76 93 b5 73 6f 9b 7f 7e 8b ec a2 13 b1 56 a6 f5 91 6e 35 83 84 9a 17 ff`.
- Step 4: From this the following HMAC key is constructed: `6f 00 11 21 3d 31 c2 3b c3 69 ab 0b 6d 8e 42 35 | 86 3a fd 35 1e 70 d5 07 76 93 b5 73 6f 9b 7f 7e 8b ec a2 13 b1 56 a6 f5 91 6e 35 83 84 9a 17 ff | 30 2d 15 d7 37 d5 b1 df 45 ee 30 bc e0 0b 89 cc` (the `|` characters illustrate the boundaries of the individual parts and do not belong to the byte values).
- Steps 5 and 6: This is the key used to calculate the SHA-3-512 HMAC value of the context id (`c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67`) which yields `8c 25 5a 6c 5a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 ad 02 d1 0f 9a 8d ae 22 6d 23 14 07 5e bc 81 c7 d3 eb 4c 71 a8 92 e7 c9 a5 6a 86 82 e4 fe f9 e7`.
- Step 7: The context key is now formed from the first 32 bytes of the HMAC value, the extended context ID and the last 32 bytes of the HMAC value: `8c 25 5a 6c 5a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 | c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67 0d | ad 02 d1 0f 9a 8d ae 22 6d 23 14 07 5e bc 81 c7 d3 eb 4c 71 a8 92 e7 c9 a5 6a 86 82 e4 fe f9 e7` (the `|` characters illustrate the boundaries of the individual parts and do not belong to the byte values).

The context id `Überführung` is thus transformed into the context key `8c 25 5a 6c 5a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67 0d ad 02 d1 0f 9a 8d ae 22 6d 23 14 07 5e bc 81 c7 d3 eb 4c 71 a8 92 e7 c9 a5 6a 86 82 e4 fe f9 e7`.

For the following hash calculations the byte sequence `8c 25 5a 6c 5a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 c3 9c 62 65 72 66 c3` is always fed in at the beginning and the byte sequence `bc 68 72 75 6a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 c3 9c 62 65 72 66 c3` is fed in at the end.

## Hash values of the files

The hash values of the files are calculated by passing the following values to the hash algorithm [SHA-3-512](https://de.wikipedia.org/wiki/SHA-3) in the following order:

1. First half of the context key
2. Bytes of the file content
3. Length of the file with variable length
4. Second half of the context key

The hash value is then read out and used for the file signature.

## Hash value of the signature file

The following procedure is used for the hash value of the signature file:

- A counter is initialized with the value `0`.
- Before each value, the counter is incremented by 1 and its value is passed in variable length.
- This is followed by the value itself.
- This is followed by the length of the value in variable length.

The hash value therefore depends on the position of a value.

The byte values are always used and not the coding.
So if a value in the signature file is encoded with Base32, the bytes that it represents are used and not the encoded values.

The values are fed in in the following order:

1. First half of the context key
2. Format identifier as a binary value, i.e. `01` for the format with the identifier `1
3. Context ID
4. Byte values of the public key
5. timestamp text
6. Computer name
7. Signature type as a binary value, i.e. `01` for `Ed25519` and `02` for ECDSAP521
8. The file names are sorted alphabetically and then fed in as follows:
    1. UTF-8 encoded name of the file
    2. Byte values of the file signature
9. Second half of the context key

The hash value is then taken from these values.

### Example

To illustrate this, the slightly modified example from the documentation is used here to show which bytes are fed into the hash algorithm:

This is the content of the example file:

```
{
   "format": 1,
   "contextId": "Überführung",
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

For the example, the value of the context ID should be `Überführung`.
This is German and means `transfer`.
This word is used here to demonstrate the encoding of non-ASCII characters in UTF-8.

The context bytes are then `c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67` and have the length 13 (`0d`).

The following values are then passed to the hash algorithm:

| Bytes                                                                                                                                                                                                                                                            | Meaning                             |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------|
| `8c 25 5a 6c 5a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 c3 9c 62 65 72 66 c3`                                                                                                                                           | 1. half of the context key          |
| `01`                                                                                                                                                                                                                                                             | Counter                             |
| `01`                                                                                                                                                                                                                                                             | Format ID                           |
| `01`                                                                                                                                                                                                                                                             | Length of format ID                 |
| `02`                                                                                                                                                                                                                                                             | Counter                             |
| `c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67`                                                                                                                                                                                                                         | Context ID                          |
| `0d`                                                                                                                                                                                                                                                             | Length of context ID                |
| `03`                                                                                                                                                                                                                                                             | Counter                             |
| `5f e2 c8 f3 98 7d 2d 5e dd df 60 87 4b f1 fc 82 13 2c d9 83 62 cc c5 37 a4 ff f0 00 ff 0b 33 86`                                                                                                                                                                | Public key                          |
| `20`                                                                                                                                                                                                                                                             | Length of context ID                |
| `04`                                                                                                                                                                                                                                                             | Counter                             |
| `32 30 32 34 2d 30 32 2d 32 35 20 31 33 3a 33 37 3a 32 32 20 2b 30 35 3a 33 30`                                                                                                                                                                                  | Timestamp                           |
| `1a`                                                                                                                                                                                                                                                             | Length of timestamp                 |
| `05`                                                                                                                                                                                                                                                             | Counter                             |
| `42 75 69 6c 64 48 6f 73 74`                                                                                                                                                                                                                                     | Hostname                            |
| `09`                                                                                                                                                                                                                                                             | Length of Hostname                  |
| `06`                                                                                                                                                                                                                                                             | Counter                             |
| `01`                                                                                                                                                                                                                                                             | Signature type                      |
| `01`                                                                                                                                                                                                                                                             | Length of signature type            |
| `07`                                                                                                                                                                                                                                                             | Counter                             |
| `63 6f 6d 6d 6f 6e 2e 67 6f`                                                                                                                                                                                                                                     | 1. file name                        |
| `09`                                                                                                                                                                                                                                                             | Length of 1. file name              |
| `08`                                                                                                                                                                                                                                                             | Counter                             |
| `ce, b2, fe, 7e, 5f, dd, bc, ec, f8, 62, b6, 49, 77, 35, bd, 36, a4, 26, a6, 18, cb, 39, 5d, ac, e7, 58, b6, 13, f5, f6, fc, 58, 1e, d3, 00, de, a9, 27, 6e, 3c, 08, 4b, 18, 39, 8c, 14, c2, 87, 1a, 50, 09, e3, eb, 31, f8, 65, 64, 00, c1, d2, cd, dc, f9, 02` | Signature of 1. file                |
| `40`                                                                                                                                                                                                                                                             | Length of signature of 1. file name |
| `09`                                                                                                                                                                                                                                                             | Counter                             |
| `66 69 6c 65 73 69 67 6e 65 72`                                                                                                                                                                                                                                  | 2. file name                        |
| `0a`                                                                                                                                                                                                                                                             | Length of 2. file name              |
| `0a`                                                                                                                                                                                                                                                             | Counter                             |
| `94, 9f, 9c, 10, c0, 03, f9, ba, 79, c0, c1, 47, af, 03, f5, 0a, fb, 8d, 74, e5, 58, 29, 37, ee, 3f, 39, 65, 95, 38, a5, a0, 06, b0, 19, cf, a6, 10, 87, f4, ed, 71, bf, a0, 53, 62, 24, 0f, ae, 3f, 5b, 54, 78, a8, 22, 31, 2d, 9e, cf, 11, e1, e0, 69, ad, 06` | Signature of 2. file                |
| `40`                                                                                                                                                                                                                                                             | Length of signature of 2. file name |
| `0b`                                                                                                                                                                                                                                                             | Counter                             |
| `66 69 6c 65 73 69 67 6e 65 72 2e 65 78 65`                                                                                                                                                                                                                      | 3. file name                        |
| `0e`                                                                                                                                                                                                                                                             | Length of 3. file name              |
| `0c`                                                                                                                                                                                                                                                             | Counter                             |
| `0c, 8f, c6, 06, 4a, f9, 33, 14, a6, 78, 9b, f2, de, b0, 58, f1, e7, 77, 97, 01, 8a, b3, 7a, 43, 00, d2, 44, 06, 59, 5d, 70, 2b, cb, be, 76, 24, 46, f7, f7, ce, df, f3, 4f, 96, 56, 31, 10, 5e, 87, a2, 59, f7, 2d, bb, c8, 6d, 3b, ab, 7f, ec, e8, 5a, 51, 08` | Signature of 3. file                |
| `40`                                                                                                                                                                                                                                                             | Length of signature of 3. file name |
| `0d`                                                                                                                                                                                                                                                             | Counter                             |
| `6d 61 70 68 65 6c 70 65 72 2f 6d 61 70 5f 68 65 6c 70 65 72 6e 65 72`                                                                                                                                                                                           | 4. file name                        |
| `17`                                                                                                                                                                                                                                                             | Length of 4. file name              |
| `0e`                                                                                                                                                                                                                                                             | Counter                             |
| `5c, f6, c8, 6a, 2b, 21, e7, af, db, 47, a9, d3, a9, 36, 6f, fa, 47, c6, 39, ca, f2, 50, 05, d3, 47, ba, 8e, 53, d4, 49, 81, 93, 6e, 92, aa, 16, 5d, b7, ff, 52, 3f, c9, 03, c2, 1d, 94, ec, a4, 8f, 9a, 5c, 8c, 1b, 21, 4f, e0, 2e, ea, ad, ac, 22, 82, 26, 0c` | Signature of 4. file                |
| `40`                                                                                                                                                                                                                                                             | Length of signature of 4. file name |
| `0f`                                                                                                                                                                                                                                                             | Counter                             |
| `73 65 74 2f 73 65 74 6e 65 72`                                                                                                                                                                                                                                  | 5. file name                        |
| `0a`                                                                                                                                                                                                                                                             | Length of 5. file name              |
| `10`                                                                                                                                                                                                                                                             | Counter                             |
| `af, 7c, 30, 97, 9a, 66, c2, f5, ae, db, fa, 64, 66, 18, 71, 47, 71, 4b, 27, 9e, 85, cd, 65, 8f, ce, b8, 08, 4a, 90, 71, 34, f3, 8d, a9, 8a, 3e, 98, 73, 7f, 27, eb, 55, 54, e8, 18, d7, 09, bf, 0a, 9d, 11, 14, 7e, ba, 63, 61, 92, 2a, 40, fd, 50, 5a, 68, 06` | Signature of 5. file                |
| `40`                                                                                                                                                                                                                                                             | Length of signature of 5. file name |
| `bc 68 72 75 6e 67 0d ad 02 d1 0f 9a 8d ae 22 6d 23 14 07 5e bc 81 c7 d3 eb 4c 71 a8 92 e7 c9 a5 6a 86 82 e4 fe f9 e7`                                                                                                                                           | 2. half of context key              |

The SHA-3-512 hash value is then generated, which is used to sign the signature file.

## Signature generation

The hash values are required to generate the signature.

In the ECDSAP521 procedure, they are used directly as described in the [Hash values of files](#hash-values-of-files) section.

With the Ed25519 procedure, however, this is much more complex.
Ed25519 requires the complete data to be signed and not just its hash values, as it calculates two hash values itself and has to go over the entire data twice to do so.
However, this is not possible here, as the files can be of any size and Ed25519 does not support a stream interface.

[RFC8032](https://www.rfc-editor.org/rfc/rfc8032) describes a variant 'Ed25519ph', where 'ph' stands for pre-hashed.
This procedure only expects a hash value.
A text with a length of 32 bytes is placed in front of the hash value and the signature of this data is calculated.
`Ed25519ph` is **not** used here.

Instead, the following procedure is used:

- `Ed25519` is used.
- Two constants are added to the left and right of the hash value.
- The signature of this extended hash value is then calculated.

The two constants have the following values:

| location | constant                                          |
|:---------|:--------------------------------------------------|
| Left     | `44 97 72 da b6 a9 2b 43 c5 06 c4 92 06 37 58 e4` |
| Right    | `b8 16 17 05 8d 38 c4 50 2b 01 2f f9 49 9e 2d dc` |

This procedure is similar to `Ed25519ph`, except that the 32 additional bytes are divided into 16 bytes each to the left and right of the hash value and other constants are used.

Example:

The hash value `ea f8 3a 32 32 e6 d0 68 ed 42 cb cf c4 7b b5 4b 28 3e c3 b6 66 54 cc c0 4e 4b 07 14 dd 02 f2 b9 58 e5 9b 05 20 aa c3 bb b5 7f d3 10 ac f9 e9 ab 5a ff 56 fa 20 5e 44 26 a0 1c 0c 3d 2a 4a ef 77` is to be signed.

The signature of the following data is then calculated using this procedure:

`44 97 72 da b6 a9 2b 43 c5 06 c4 92 06 37 58 e4 ea f8 3a 32 32 e6 d0 68 ed 42 cb cf c4 7b b5 4b 28 3e c3 b6 66 54 cc c0 4e 4b 07 14 dd 02 f2 b9 58 e5 9b 05 20 aa c3 bb b5 7f d3 10 ac f9 e9 ab 5a ff 56 fa 20 5e 44 26 a0 1c 0c 3d 2a 4a ef 77 b8 16 17 05 8d 38 c4 50 2b 01 2f f9 49 9e 2d dc`
