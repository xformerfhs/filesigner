# Technische Spezifikation

## Allgemein

In dieser Datei werden die technischen Details der Berechnung und Prüfung der Signaturen beschrieben.
Im ersten Teil werden die kryptografischen Grundlagen erläutert.
Im zweiten Teil werden die einzelnen Berechnungsschritte spezifiziert.

## Kryptografische Grundlagen

Dieses System arbeitet mit digitalen Signaturen.
Folgende Schritte werden zur Erstellung einer digitalen Signatur benötigt:

1. Ein kryptografisch sicheres [Hash-Verfahrens](https://de.wikipedia.org/wiki/Kryptographische_Hashfunktion).
2. Ein [asymmetrischen Verschlüsselungsverfahrens](https://de.wikipedia.org/wiki/Asymmetrisches_Kryptosystem) wie [RSA](https://de.wikipedia.org/wiki/RSA-Kryptosystem) oder [elliptische Kurven](https://de.wikipedia.org/wiki/Elliptic_Curve_Cryptography).

### Signatur

Für die Erstellung einer Signatur wird zuerst mit einem krytopraphisch sicheren Hash-Verfahren eine kryptographisch sichere [Prüfsumme](https://de.wikipedia.org/wiki/Pr%C3%BCfsumme) über die zu signierenden Daten ermittelt.
Dieser kryptographisch sichere Hash-Wert wird dann mit dem **privaten** - also geheimen - Schlüssel eines asymmetrischen Verschlüsselungsverfahren verschlüsselt.
Dieser so verschlüsselte Hash_Wert ist die digitale Signatur.

Bei einem asymmetrischen Verschlüsselungsverfahren gibt es ein Schlüsselpaar aus einem privaten und einem öffentlichen Schlüssel.
Dabei gelten die folgenden Regeln bzgl. der Ver- und Entschlüsselung:

- Was mit dem **öffentlichen** Schlüssel __verschlüsselt__ wird, kann nur mit dem **privaten** Schlüssel __entschlüsselt__ werden.
- Was mit dem **privaten** Schlüssel __verschlüsselt__ wird, kann nur mit dem **öffentlichen** Schlüssel __entschlüsselt__ werden.

D.h., dass alles, was mit dem **einen** Schlüssel verschlüsselt wird, immer nur mit dem **anderen** Schlüssel des Schlüsselpaars entschlüsselt werden kann.

Wie der Name bereits sagt, ist der öffentliche Schlüssel bekannt und wird nicht geheim gehalten.
Der private Schlüssel muss geheim gehalten werden.

Für die Signatur wird der ermittelte Hash-Wert mit dem **privaten** - also geheimen - Schlüssel des Erstellers verschlüsselt.
Der so verschlüsselte Hash-Wert stellt die digitale Signatur dar.

### Verifizierung

For verification the verifier calculates the hash value of the received data using the same hash procedure as the creator of the message.
Then the signature is decrypted with the **public** - i.e. known - key of the creator.
The decrypted hash value from the signature is compared to the hash value that has been calculated locally.

Zur Verifizierung berechnet der Verifizierer als erstes den Hash-Wert der empfangenen Daten mit demselben Hash-Verfahren, dass auch der Unterschreibende der Daten benutzt hat.
Dann entschlüsselt er die digitale Signatur mit dem **öffentlichen** - also bekannten - Schlüssel des Unterschreibenden.
Der entschlüsselte Hash-Wert mit dann mit dem lokal erzeugten Hash-Wert verglichen.
Stimmen die beiden Hash-Werte überein, ist die digitale Signatur gültig, andernfalls nicht.

### Bedeutung

Wenn die Verifizierung gelingt, ist sichergestellt, dass der Ersteller der Signatur den zum öffentlichen Schlüssel passenden privaten Schlüssel benutzt hat.

Die Sicherheit dieses Verfahrens ruht auf zwei Säulen:

1. Das Hash-Verfahren stellt sicher, dass ein Angreifer die Daten nicht so verändern kann, dass mit den veränderten Daten derselbe Hash-Wert entsteht.
2. Das Verschlüsselungsverfahren stellt sicher, dass aus der Kenntnis des öffentlichen Schlüssels nicht auf den Wert des privaten Schlüssels geschlossen werden kann.

### Ausgewählte Verfahren

Für die Berechnung des Hash-Wertes wird das Verfahren SHA-3-512 benutzt, also [SHA-3](https://de.wikipedia.org/wiki/SHA-3) mit einer Hash-Länge von 512 Bit (64 Byte).
Dieses Verfahren wurde vom [NIST](https://www.nist.gov/) standardisiert und ist das zurzeit sicherste Hash-Verfahren mit einer sehr langen und damit noch auf lange Sicht sicheren Länge des Hash-Wertes.

Als Signaturverfahren werden [Ed25519](https://de.wikipedia.org/wiki/Curve25519#Ed25519_und_weitere_Kurven) und [ECDSA](https://de.wikipedia.org/wiki/Elliptic_Curve_DSA) mit der Kurve [P-521](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-186.pdf) verwendet.

Beide Verfahren benutzen elliptische Kurven als asymmetrisches Verschlüsselungsverfahren.
Elliptische Kurven sind zur Zeit und auf absehbare Zeit sicher gegen Angriffe durch klassische Computer.
Theoretisch sind sie durch Quantencomputer angreifbar.
Sie sind aber auch bei Quantencomputern sicherer, als das RSA-Verfahren, da ein Angriff auf einem Quantencomputer für elliptische Kurven komplexer ist.
Sollte eine konkrete Gefahr für elliptische Kurven bekannt werden, kann das System jederzeit so geändert werden, dass ein quantencomputersicheres Signatursystem wie [CRYSTALS-Dilithium](https://pq-crystals.org/dilithium/), [FALCON](https://falcon-sign.info/) oder [SPHINCS+](https://sphincs.org/) verwendet wird.
Zurzeit sind aber kaum Bibliotheken für die Nutzung dieser Systeme verfügbar.

Gleichzeitig sind elliptische Kurven wegen ihrer gegenüber RSA wesentlich kürzeren Schlüssellängen effizienter im Ressourcenverbrauch.
Zusammengefasst sind zur Zeit keine effektiven Angriffe auf elliptische Kurven bekannt.

Das Verfahren Ed25519 basiert auf der elliptischen Kurve [Curve25519](https://www.rfc-editor.org/rfc/rfc8032#page-9), die die sicherste bekannte elliptische Kurve ist und die effizient berechenbar ist.

Das Verfahren ECDSA mit der Kurve P-521 benutzt dagegen die vom NIST standardisierte Kurve P-521 mit dem ebenfalls vom NIST erstellten Standard ECDSA.
Das Verfahren ist deutlich ineffizienter, als Ed25519, dafür aber in vielen Programmiersprachen vorhanden.

## Berechnungen

Im folgenden Abschnitt wird beschrieben, wie die einzelnen Berechnungen durchgeführt werden.

Wann immer der Inhalt von Bytes angegeben ist, werden diese in der hexadezimalen Schreibweise angegeben.

Zwei Werte werden dabei mehrmals benutzt, so dass diese hier im Vorfeld beschrieben werden:

### Zähler variabler Länge

Für die Berechnung werden Zähler oder Längen benutzt.
Deren Werte werden immer mit so vielen Bytes benutzt, wie nötig sind und nicht mehr.
Die Bytes werden im sogenannten [Big-Endian](https://de.wikipedia.org/wiki/Byte-Reihenfolge#Big-Endian-Format)-Format angegeben.

Zur Erläuterung hier ein paar Beispiele für Werte mit ihrer Kodierung in variabler Länge:

|     Wert | Kodierung (hexadezimal) |
|---------:|-------------------------|
|      `0` | `00`                    |
|      `1` | `01`                    |
|    `255` | `ff`                    |
|    `300` | `01 2c`                 |
|  `65432` | `ff 98`                 |
| `100000` | `01 86 a0`              |

### Kontext-Schlüssel

Die Signaturen benötigen eine "Kontext-Id".
Diese Kontext-Id fließt in die Hash-Berechnung jeder Datei ein.
Allerdings wird die Kontext-Id vom Aufrufer gesteuert, so dass es einem Angreifer möglich wäre durch eine manipulierte Kontext-Id die Signaturverfahren anzugreifen[^1].
Um solche Angriffe zu verhindern, wird aus der Kontext-Id ein Kontext-Schlüssel berechnet.
Dieser Kontext-Schlüssel fließt dann in die Hash-Berechnungen der Dateien ein, nicht die Kontext-Id selbst.

[^1]: Es sind zur Zeit keine solchen Angriffe bekannt.

Der Kontext-Schlüssel wird folgendermaßen aus der Kontext-Id berechnet:

1. Die Zeichen der `Kontext-Id` werden in [`UTF-8`](https://de.wikipedia.org/wiki/UTF-8) kodiert.
2. An die so entstandene Byte-Folge wird die Länge in der Darstellung mit variable Länge angehängt, was im Folgenden als "erweiterter Kontext-Id" bezeichnet wird.
3. Die Byte-Folge der erweiterten Kontext-Id wird herumgedreht und daraus der SHA-3-256-Wert berechnet.
4. Aus diesem Hash-Wert mit einer Länge von 32 Byte wird mit den folgenden Werten ein 64-Byte-langer Schlüssel konstruiert:
  - Konstante Byte-Folge `6f 00 11 21 3d 31 c2 3b c3 69 ab 0b 6d 8e 42 35`.
  - Der eben berechnete Hash-Wert.
  - Konstante Byte-Folge `30 2d 15 d7 37 d5 b1 df 45 ee 30 bc e0 0b 89 cc`.
5. Dieser Schlüssel mit einer Länge von 64 Byte dient als Schlüssel für ein [HMAC](https://de.wikipedia.org/wiki/HMAC)-SHA-3-512-Verfahren genutzt.
6. Es wird der HMAC-Wert der Kontext-Id mit dem so konstruierten Schlüssel berechnet.
7. Aus dem so erzeugte 64 byte langen HMAC-Wert wird mit den folgenden Werten ein Padding erzeugt:
  - Die ersten 32 Bytes des HMAC-Wertes.
  - Die erweiterte Kontext-Id.
  - Die letzten 32 Bytes des HMAC-Wertes.

Der so errechnete Schlüssel wird bei der Berechnung aller Hash-Werte von Dateien in der folgeden Weise benutzt.
  - Zuerst wird die erste Hälfte des Schlüssels in den Hash-Wert eingespeist.
  - Dann werden alle Daten der Datei in den Hash-Wert eingespeist.
  - Am Schluss wird die zweite Hälfte des Schlüssels in den Hash-Wert eingespeist.
  - Bei einer ungeraden Schlüssellänge ist die erste Hälfte des Schlüssels um ein Byte kürzer, als die zweite Hälfte.

Ein Beispiel zeigt diese Berechnungsvorschrift an konkreten Werten:

- Die Kontext-Id lautet `Überführung`.
- Schritt 1: Die Byte-Folge der Kontext-Id lautet in der UTF-8-Kodierung `c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67` und hat die Länge 13 (`0d`).
- Schritt 2: Die erweiterte Kontex-Id hat damit den Wert`c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67 0d`.
- Schritt 3: Der SHA-3-256-Wert der herumgedrehten erweiterten Kontext-Id (`0d 67 6e 75 72 68 bc c3 66 72 65 62 9c c3`) wird berechnet und ergibt `86 3a fd 35 1e 70 d5 07 76 93 b5 73 6f 9b 7f 7e 8b ec a2 13 b1 56 a6 f5 91 6e 35 83 84 9a 17 ff`.
- Schritt 4: Damit wird der folgende HMAC-Schlüssel konstruiert: `6f 00 11 21 3d 31 c2 3b c3 69 ab 0b 6d 8e 42 35 | 86 3a fd 35 1e 70 d5 07 76 93 b5 73 6f 9b 7f 7e 8b ec a2 13 b1 56 a6 f5 91 6e 35 83 84 9a 17 ff | 30 2d 15 d7 37 d5 b1 df 45 ee 30 bc e0 0b 89 cc` (die `|`-Zeichen zeigen die Grenzen zwischen den einzelnen Teilen der Konstruktion an und gehören nicht zu den Byte-Werten der Konstruktion).
- Steps 5 and 6: Dieser so konstruierte Schlüssel wird benutzt, um den SHA-3-512-HMAC-Wert der Kontext-Id (`c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67`) zu berechnen und ergibt `8c 25 5a 6c 5a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 ad 02 d1 0f 9a 8d ae 22 6d 23 14 07 5e bc 81 c7 d3 eb 4c 71 a8 92 e7 c9 a5 6a 86 82 e4 fe f9 e7`.
- Schritt 7:Der Kontext-Schlüssel wird nun aus den ersten 32 Byte dieses HMAC-Wertes, der erweiterten Kontext-id und den letzten 32 Byte dieses HMAC-Wertes gebildet: `8c 25 5a 6c 5a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 | c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67 0d | ad 02 d1 0f 9a 8d ae 22 6d 23 14 07 5e bc 81 c7 d3 eb 4c 71 a8 92 e7 c9 a5 6a 86 82 e4 fe f9 e7` (die `|`-Zeichen zeigen die Grenzen zwischen den einzelnen Teilen der Konstruktion an und gehören nicht zu den Byte-Werten der Konstruktion).

Damit ist die Kontext-Id `Überführung` in den Kontext-Schlüssel `8c 25 5a 6c 5a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67 0d ad 02 d1 0f 9a 8d ae 22 6d 23 14 07 5e bc 81 c7 d3 eb 4c 71 a8 92 e7 c9 a5 6a 86 82 e4 fe f9 e7` überführt worden.

In den folgenden Hash-Berechnungen wird am Anfang immer die Bytefolge `8c 25 5a 6c 5a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 c3 9c 62 65 72 66 c3` und am Ende immer die Bytefolge `bc 68 72 75 6a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 c3 9c 62 65 72 66 c3` eingespeist.

## Hash-Werte der Dateien

Die Hash-Werte der Dateien werden berechnet, indem die folgenden Werte in folgender Reihenfolge an den Hash-Algorithmus [SHA-3-512](https://de.wikipedia.org/wiki/SHA-3) übergeben werden:

1. Erste Hälfte des Kontext-Schlüssels
2. Bytes des Dateiinhaltes
3. Länge der Datei mit variabler Länge
4. Zweite Hälfte des Kontext_Schlüssels

Danach wird der Hash-Wert ausgelesen und für die Dateisignatur benutzt.

## Hash-Wert der Signaturendatei

Für den Hash-Wert der Signaturendatei wird folgendes Verfahren angewendet:

- Es wird ein Zähler mit dem Wert `0` initialisiert.
- Vor jedem Wert wird der Zähler um 1 erhöht und sein Wert in variabler Länge übergeben.
- Dann folgt der Wert selbst.
- Dann folgt die Länge des Wertes in variabler Länge.

Damit hängt der Hash-Wert von der Position eines Wertes ab.

Es werden immer die Byte-Werte benutzt und nicht die Kodierung.
Wenn also ein Wert in der Signaturendatei mit Base32 kodiert ist, werden die Bytes benutzt, die er repräsentiert und nicht die kodierten Werte.

Die Werte werden in der folgenden Reihenfolge eingespeist:

1. Erste Hälfte des Kontext-Schlüssels
2. Die Format-Kennung als Binärwert, also `01` für das Format mit der Kennung `1`.
3. Die Kontext-Id
4. Die Byte-Werte des öffentlichen Schlüssels
5. Der Text des Zeitstempels
6. Der Text des Rechnernamens
7. Der Signaturtyp als Binärwert, also `01` für `Ed25519` und `02` für ECDSAP521
8. Die Dateinamen werden alphabetisch sortiert und dann jeweils folgendermaßen eingespeist:
    1. Der Name der Datei in UTF-8-Kodierung
    2. Die Byte-Werte der Signatur der Datei
9. Die zweite Hälfte des Kontext-Schlüssels

Danach wird der Hash-Wert aus diesen Werten entnommen.

### Beispiel

Zur Illustration wird hier das etwas abgewandelte Beispiel aus der Dokumentation genommen und gezeigt, welche Bytes in den Hash-Algorithmus gespeist werden:

Dies ist der Inhalt der Beispiel-Datei:

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

Für das Beispiel soll der Wert der Kontext-Id `Überführung` sein.
Die Kontext-Bytes lauten dann `c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67` und haben die Länge 13 (`0d`).

Folgende Werte werden dann an den Hash-Algorithmus übergeben:

| Bytes                                                                                                                                                                                                                                                            | Bedeutung                             |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------|
| `8c 25 5a 6c 5a 75 d2 ab bc 34 c7 2f 38 a8 da db 7b 39 97 47 b1 9e 3e e8 d3 9a f9 cf 83 9a 39 03 c3 9c 62 65 72 66 c3`                                                                                                                                           | 1. Hälfte des Kontext-Schlüssels      |
| `01`                                                                                                                                                                                                                                                             | Zähler                                |
| `01`                                                                                                                                                                                                                                                             | Format-Kennung                        |
| `01`                                                                                                                                                                                                                                                             | Länge der Format-Kennung              |
| `02`                                                                                                                                                                                                                                                             | Zähler                                |
| `c3 9c 62 65 72 66 c3 bc 68 72 75 6e 67`                                                                                                                                                                                                                         | Kontext-Id                            |
| `0d`                                                                                                                                                                                                                                                             | Länge der Kontext-Id                  |
| `03`                                                                                                                                                                                                                                                             | Zähler                                |
| `5f e2 c8 f3 98 7d 2d 5e dd df 60 87 4b f1 fc 82 13 2c d9 83 62 cc c5 37 a4 ff f0 00 ff 0b 33 86`                                                                                                                                                                | Öffentlicher Schlüssel                |
| `20`                                                                                                                                                                                                                                                             | Länge des öffentlichen Schlüssels     |
| `04`                                                                                                                                                                                                                                                             | Zähler                                |
| `32 30 32 34 2d 30 32 2d 32 35 20 31 33 3a 33 37 3a 32 32 20 2b 30 35 3a 33 30`                                                                                                                                                                                  | Text des Zeitstempels                 |
| `1a`                                                                                                                                                                                                                                                             | Länge des Zeitstempels                |
| `05`                                                                                                                                                                                                                                                             | Zähler                                |
| `42 75 69 6c 64 48 6f 73 74`                                                                                                                                                                                                                                     | Host-Name                             |
| `09`                                                                                                                                                                                                                                                             | Länge des Host-Namens                 |
| `06`                                                                                                                                                                                                                                                             | Zähler                                |
| `01`                                                                                                                                                                                                                                                             | Signaturtyp                           |
| `01`                                                                                                                                                                                                                                                             | Länge des Signaturtyps                |
| `07`                                                                                                                                                                                                                                                             | Zähler                                |
| `63 6f 6d 6d 6f 6e 2e 67 6f`                                                                                                                                                                                                                                     | 1. Dateiname                          |
| `09`                                                                                                                                                                                                                                                             | Länge des 1. Dateinamens              |
| `08`                                                                                                                                                                                                                                                             | Zähler                                |
| `ce, b2, fe, 7e, 5f, dd, bc, ec, f8, 62, b6, 49, 77, 35, bd, 36, a4, 26, a6, 18, cb, 39, 5d, ac, e7, 58, b6, 13, f5, f6, fc, 58, 1e, d3, 00, de, a9, 27, 6e, 3c, 08, 4b, 18, 39, 8c, 14, c2, 87, 1a, 50, 09, e3, eb, 31, f8, 65, 64, 00, c1, d2, cd, dc, f9, 02` | Signatur der 1. Datei                 |
| `40`                                                                                                                                                                                                                                                             | Länge des Signatur des 1. Dateinamens |
| `09`                                                                                                                                                                                                                                                             | Zähler                                |
| `66 69 6c 65 73 69 67 6e 65 72`                                                                                                                                                                                                                                  | 2. Dateiname                          |
| `0a`                                                                                                                                                                                                                                                             | Länge des 2. Dateinamens              |
| `0a`                                                                                                                                                                                                                                                             | Zähler                                |
| `94, 9f, 9c, 10, c0, 03, f9, ba, 79, c0, c1, 47, af, 03, f5, 0a, fb, 8d, 74, e5, 58, 29, 37, ee, 3f, 39, 65, 95, 38, a5, a0, 06, b0, 19, cf, a6, 10, 87, f4, ed, 71, bf, a0, 53, 62, 24, 0f, ae, 3f, 5b, 54, 78, a8, 22, 31, 2d, 9e, cf, 11, e1, e0, 69, ad, 06` | Signatur der 2. Datei                 |
| `40`                                                                                                                                                                                                                                                             | Länge des Signatur des 2. Dateinamens |
| `0b`                                                                                                                                                                                                                                                             | Zähler                                |
| `66 69 6c 65 73 69 67 6e 65 72 2e 65 78 65`                                                                                                                                                                                                                      | 3. Dateiname                          |
| `0e`                                                                                                                                                                                                                                                             | Länge des 3. Dateinamens              |
| `0c`                                                                                                                                                                                                                                                             | Zähler                                |
| `0c, 8f, c6, 06, 4a, f9, 33, 14, a6, 78, 9b, f2, de, b0, 58, f1, e7, 77, 97, 01, 8a, b3, 7a, 43, 00, d2, 44, 06, 59, 5d, 70, 2b, cb, be, 76, 24, 46, f7, f7, ce, df, f3, 4f, 96, 56, 31, 10, 5e, 87, a2, 59, f7, 2d, bb, c8, 6d, 3b, ab, 7f, ec, e8, 5a, 51, 08` | Signatur der 3. Datei                 |
| `40`                                                                                                                                                                                                                                                             | Länge des Signatur des 3. Dateinamens |
| `0d`                                                                                                                                                                                                                                                             | Zähler                                |
| `6d 61 70 68 65 6c 70 65 72 2f 6d 61 70 5f 68 65 6c 70 65 72 6e 65 72`                                                                                                                                                                                           | 4. Dateiname                          |
| `17`                                                                                                                                                                                                                                                             | Länge des 4. Dateinamens              |
| `0e`                                                                                                                                                                                                                                                             | Zähler                                |
| `5c, f6, c8, 6a, 2b, 21, e7, af, db, 47, a9, d3, a9, 36, 6f, fa, 47, c6, 39, ca, f2, 50, 05, d3, 47, ba, 8e, 53, d4, 49, 81, 93, 6e, 92, aa, 16, 5d, b7, ff, 52, 3f, c9, 03, c2, 1d, 94, ec, a4, 8f, 9a, 5c, 8c, 1b, 21, 4f, e0, 2e, ea, ad, ac, 22, 82, 26, 0c` | Signatur der 4. Datei                 |
| `40`                                                                                                                                                                                                                                                             | Länge des Signatur des 4. Dateinamens |
| `0f`                                                                                                                                                                                                                                                             | Zähler                                |
| `73 65 74 2f 73 65 74 6e 65 72`                                                                                                                                                                                                                                  | 5. Dateiname                          |
| `0a`                                                                                                                                                                                                                                                             | Länge des 5. Dateinamens              |
| `10`                                                                                                                                                                                                                                                             | Zähler                                |
| `af, 7c, 30, 97, 9a, 66, c2, f5, ae, db, fa, 64, 66, 18, 71, 47, 71, 4b, 27, 9e, 85, cd, 65, 8f, ce, b8, 08, 4a, 90, 71, 34, f3, 8d, a9, 8a, 3e, 98, 73, 7f, 27, eb, 55, 54, e8, 18, d7, 09, bf, 0a, 9d, 11, 14, 7e, ba, 63, 61, 92, 2a, 40, fd, 50, 5a, 68, 06` | Signatur der 5. Datei                 |
| `40`                                                                                                                                                                                                                                                             | Länge des Signatur des 5. Dateinamens |
| `bc 68 72 75 6e 67 0d ad 02 d1 0f 9a 8d ae 22 6d 23 14 07 5e bc 81 c7 d3 eb 4c 71 a8 92 e7 c9 a5 6a 86 82 e4 fe f9 e7`                                                                                                                                           | 2. Hälfte des Kontext-Schlüssels      |

Danach wird daraus der SHA-3-512-Hash-Wert erzeugt, der für die Signatur der Signaturen-Datei benutzt wird.

## Signaturerzeugung

Die Hash-Werte werden für die Erzeugung der Signatur benötigt.

Beim Verfahren ECDSAP521 werden sie direkt so verwendet, wie sie von im Abschnitt [Hash-Werte der Dateien](#hash-werte-der-dateien) beschrieben wurden.

Beim Verfahren Ed25519 ist das jedoch deutlich komplexer.
Ed25519 benötigt die vollständigen zu signierenden Daten und nicht nur deren Hash-Werte, da es selbst zwei Hash-Werte berechnet und dafür zweimal über die ganzen Daten gehen muss.
Das ist jedoch hier nicht möglich, da die Dateien beliebig groß sein können und Ed25519 kein Stream-Interface unterstützt.

Im [RFC8032](https://www.rfc-editor.org/rfc/rfc8032) wird eine Variante `Ed25519ph` beschrieben, wobei `ph` für pre-hashed steht.
Dieses Verfahren erwartet nur einen Hash-Wert.
Ein Text mit 32 Bytes Länge wird dem Hash-Wert vorangestellt und die Signatur dieser Daten berechnet. 
`Ed25519ph` wird hier **nicht** benutzt.

Stattdessen wird folgendes Verfahren verwendet:

- Es wird `Ed25519` benutzt.
- Der Hash-Wert wird links und rechts mit zwei Konstanten ergänzt.
- Es wird dann die Signatur dieses erweiterten Hash-Wertes berechnet.

Die beiden Konstanten haben die folgenden Werte:

| Ort    | Konstante                                           |
|:-------|:----------------------------------------------------|
| Links  | `44 97 72 da b6 a9 2b 43 c5 06 c4 92 06 37 58 e4`   |
| Rechts | `b8 16 17 05 8d 38 c4 50 2b 01 2f f9 49 9e 2d dc`   |

Dieses Verfahren ähnelt `Ed25519ph`, nur dass die 32 zusätzlichen Bytes auf je 16 Bytes links und rechts des Hash-Wertes aufgeteilt und andere Konstanten benutzt werden.

Beispiel:

Der Hash-Wert `ea f8 3a 32 32 e6 d0 68 ed 42 cb cf c4 7b b5 4b 28 3e c3 b6 66 54 cc c0 4e 4b 07 14 dd 02 f2 b9 58 e5 9b 05 20 aa c3 bb b5 7f d3 10 ac f9 e9 ab 5a ff 56 fa 20 5e 44 26 a0 1c 0c 3d 2a 4a ef 77` soll signiert werden.

Dann wird mit diesem Verfahren die Signatur der folgenden Daten berechnet:

`44 97 72 da b6 a9 2b 43 c5 06 c4 92 06 37 58 e4 ea f8 3a 32 32 e6 d0 68 ed 42 cb cf c4 7b b5 4b 28 3e c3 b6 66 54 cc c0 4e 4b 07 14 dd 02 f2 b9 58 e5 9b 05 20 aa c3 bb b5 7f d3 10 ac f9 e9 ab 5a ff 56 fa 20 5e 44 26 a0 1c 0c 3d 2a 4a ef 77 b8 16 17 05 8d 38 c4 50 2b 01 2f f9 49 9e 2d dc`
