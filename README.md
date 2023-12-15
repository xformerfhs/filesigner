# filesigner

A program to have an easy way to use digital signatures.
Like certificates, only better and easier ;-).

## Introduction

We create artifacts and store them in different systems.
Currently there is no way to check if anything has been changed in the stored artifacts.
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

### Signature

The signature call looks like this:

```
filesigner sign {contextId} [type] ! {fileList}
```

Where `contextId` is an arbitrary text used to make the signature depend on a topic.
For example, one could use the GitLab pipeline number or some other attribute that matches the creation of the artifact.

After the context id, there may or may not be a specification of the signature method.
This is either `ed25519` or `ecdsap521`.
If the type is not specified, `ed25519` is used.

This is followed by a single exclamation mark (`!`).
This is used to separate the `contextId` and possibly the signature type logically and visually from the file list.

This is followed by a list of file names that are to be signed.
Wildcards like `*` and `?` can be used.

If a filename starts with a `-`, a file with the following name, or with a name matching the following pattern, will **not** be signed.
E.g. a specification of `-*.exe` excludes all files with the extension `.exe`.
The file `signatures.json` is **always** excluded and cannot be signed.
It already contains a signature.
On Linux, wildcard specifications for exclusion must always be enclosed in apostrophes (`'`).

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
| `0`  | Successful processing.     |
| `1`  | Error in the command line.   |
| `2`  | Warning while processing. |
| `3`  | Error while processing.  |

### Verifizierung

Der Aufruf zur Verifizierung sieht folgendermaßen aus:

```
filesigner verify {contextId} 
```

Dabei hat `contextId` dieselbe Bedeutung, wie bei der Signatur.

Weitere Parameter sind nicht erlaubt und führen zu einer Fehlermeldung.

Das Programm liest die Datei `signatures.json` ein und prüft, ob die dort genannten Dateien vorhanden sind und ob deren Signaturen zu den aktuellen Inhalten passen.

Die Rückgabe-Codes sind dieselben, wie bei der Signierung.

## Programme

| OS      | Programm             |
|---------|----------------------|
| Windows | `filesigner.exe` |
| Linux   | `filesigner`     |

Das Linux-Programm ist auf jedem beliebigen Linux ausführbar.

## Kodierung

Binärwerte sind in einer speziellen [Base32-Kodierung](https://en.wikipedia.org/wiki/Base32) abgelegt.
Die Besonderheit besteht darin, dass das verwendete Alphabet keine Vokale, keine leicht zu verwechselnden Zeichen wie `0` und `O` oder `1` und `l` und keine Sonderzeichen enthält.
Dadurch können die so kodierten Werte mit einem Doppelklick markiert werden, es können nicht versehentlich echte Worte entstehen und beim Vorlesen kann man keine Zeichen verwechseln.

## Beispiel

Angenommen, für die Artefakte `filesigner`, `filesigner.exe`, alle `Go`-Dateien und alle Dateien, die mit dem Wort `go`beginnen, sollen für die Version `1.7.11` einer Anwendung Signaturen erzeugt und überprüft werden.

Dann erzeugt man die Signatur mit folgendem Aufruf:

```
filesigner sign projekt1711 ! *.go filesign*
```

Das Programm erzeugt dann die folgende Ausgabe auf der Konsole:

```
2023-06-13 12:51:25 +02:00  18  I  filesigner V0.13.4 (go1.21.3, 12 cpus)
2023-06-13 12:51:25 +02:00  39  I  Context id         : charm
2023-06-13 12:51:25 +02:00  40  I  Public key id      : 85R3-VZPX-JRV8-RN6R-G0L1-SV4U-NW
2023-06-13 12:51:25 +02:00  41  I  Signature timestamp: 2023-06-13 12:51:25 +02:00
2023-06-13 12:51:25 +02:00  42  I  Signature host name: MDXN01022044
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'common.go'
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'filesigner'
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'filesigner.exe'
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'main.go'
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'sign_command.go'
2023-06-13 12:51:25 +02:00  21  I  Signing succeeded for file 'verify_command.go'
2023-06-13 12:51:25 +02:00  43  I  Signatures for 6 files successfully created
```

Der Rückgabe-Code ist 0.

Die Artefakte werden zusammen mit der Datei in einem Repository abgelegt und das Erzeuger-Team informiert die Empfänger mit einer signierten E-Mail, dass die Artefakte am Zielort liegen.
In dieser E-Mail steht, dass die Signatur mit dem öffentlichen Schlüssel mit der Schlüssel-Id `85R3-VZPX-JRV8-RN6R-G0L1-SV4U-NW` und dem Kontext `projekt1711` auf dem Host `MDXN01022044` zum Zeitpunkt `2023-06-13 12:51:25 +02:00` erzeugt wurde und geprüft werden kann.

Die Empfänger laden die Dateien herunter und prüfen, ob in der Datei `signatures.json` der angegebene öffentliche Schlüssel steht.
Dann rufen sie das Programm folgendermaßen auf:

```
filesigner verify projekt1711
```

Das Programm erzeugt dann die folgende Ausgabe auf der Konsole:

```
2023-06-13 12:51:52 +02:00  18  I  filesigner V0.13.4 (go1.21.3, 12 cpus)
2023-06-13 12:51:53 +02:00  58  I  Context id         : charm
2023-06-13 12:51:53 +02:00  59  I  Public key id      : 85R3-VZPX-JRV8-RN6R-G0L1-SV4U-NW
2023-06-13 12:51:53 +02:00  60  I  Signature timestamp: 2023-06-13 12:51:25 +02:00
2023-06-13 12:51:53 +02:00  61  I  Signature host name: MDXN01022044
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'common.go'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'filesigner'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'filesigner.exe'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'main.go'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'sign_command.go'
2023-06-13 12:51:53 +02:00  21  I  Verification succeeded for file 'verify_command.go'
2023-06-13 12:51:53 +02:00  64  I  Verification of 6 files successful
```

Der Rückgabe-Code ist 0.

Damit sind die Signaturen verifiziert.

Sollte z.B. die Datei `filesigner` manipuliert worden sein, würde folgende Ausgabe erscheinen:

```
2023-06-13 12:54:05 +02:00  18  I  filesigner V0.13.4 (go1.21.3, 12 cpus)
2023-06-13 12:54:05 +02:00  58  I  Context id         : charm
2023-06-13 12:54:05 +02:00  59  I  Public key id      : 85R3-VZPX-JRV8-RN6R-G0L1-SV4U-NW
2023-06-13 12:54:05 +02:00  60  I  Signature timestamp: 2023-06-13 12:51:25 +02:00
2023-06-13 12:54:05 +02:00  61  I  Signature host name: MDXN01022044
2023-06-13 12:54:05 +02:00  21  I  Verification succeeded for file 'common.go'
2023-06-13 12:54:05 +02:00  21  I  Verification succeeded for file 'filesigner.exe'
2023-06-13 12:54:05 +02:00  21  I  Verification succeeded for file 'main.go'
2023-06-13 12:54:05 +02:00  21  I  Verification succeeded for file 'sign_command.go'
2023-06-13 12:54:05 +02:00  21  I  Verification succeeded for file 'verify_command.go'
2023-06-13 12:54:05 +02:00  22  E  File 'filesigner' has been tampered with
2023-06-13 12:54:05 +02:00  66  E  Verification of 5 files successful and 1 file unsuccessful
```

Der Rückgabe-Code ist 3.

Sollte z.B. die Signaturdatei manipuliert worden sein oder die falsche Kontext-Id angegeben worden sein, würde folgende Ausgabe erscheinen:

```
2023-06-13 12:54:56 +02:00  18  I  filesigner V0.13.4 (go1.21.3, 12 cpus)
2023-06-13 12:54:56 +02:00  58  I  Context id         : charm
2023-06-13 12:54:56 +02:00  59  I  Public key id      : 85R3-VZPX-JRV8-RN6R-G0L1-SV4U-NW
2023-06-13 12:54:56 +02:00  60  I  Signature timestamp: 2023-06-13 12:51:25 +02:00
2023-06-13 12:54:56 +02:00  61  I  Signature host name: MDXN01022044
2023-06-13 12:54:56 +02:00  62  E  Signature file has been tampered with or wrong context id
```

Der Rückgabe-Code ist 3.

## Erstellung

Zur Erstellung des Programms muss man Go installiert haben.
Dabei wird ein Verzeichnis angelegt, dass in der Umgebungsvariablen `GOPATH` spezifiziert ist.
Unter Windows ist das das Heimatverzeichnis, z.B. `D:\Users\Benutzername\go`.
Unter Linux ist es `${HOME}/go`.
Darunter befindet sich ein Verzeichnis mit dem Namen `src`.

Zum Erstellen des Programms müssen die Quelltexte unter `%GOPATH%\src\filesigner`, bzw. `${HOME}/go/src/filesigner` abgelegt werden.
Danach ruft man die Batch-Datei `gb.bat`, bzw. das Shell-Skript `gb` auf, die die Erstellung übernehmen.
Sie erwarten das Programm UPX an einem bestimmten Ort.
Diesen Ort kann man auf den lokal vorhandenen Pfad anpassen.
Wenn UPX nicht vorhanden ist, wird keine Komprimierung durchgeführt.

Als Ergebnis werden die Dateien `filesigner` für Linux und `filesigner.exe` für Windows erstellt.

Die Artefakte findet man in unserem Artifactory unter der [Web-Oberfläche](https://bahnhub.tech.rz.db.de/ui/repos/tree/General/davit-generic-stage-dev-local/zzz/tools/filesigner/) oder als [Download](https://bahnhub.tech.rz.db.de:443/artifactory/davit-generic-stage-dev-local/zzz/tools/filesigner/).

## Kontakt

Frank Schwab ([Mail](mailto:frank.schwab@deutschebahn.com "Mail"))

## Lizenz

Das Programm ist unter der [Apache-Lizenz V2](https://www.apache.org/licenses/LICENSE-2.0.txt) veröffentlicht.
