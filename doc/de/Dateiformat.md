# Dateiformat

## Beschreibung

In diesem Dokument wird der Aufbau der Signaturendatei beschrieben.

Die Daten sind im [JSON](https://de.wikipedia.org/wiki/JavaScript_Object_Notation)-Format abgelegt.
Alle Texte sind in [UTF-8](https://de.wikipedia.org/wiki/UTF-8) kodiert.

In der Datei sind die folgenden Felder vorhanden:

| Feld             | Bedeutung                                                                                                                                 |
|------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| `contextId`      | Die Kontext-Id der Signatur.                                                                                                              |
| `dataSignature`  | Die Signatur über die einzelnen Teile dieser Datei.                                                                                       |
| `fileSignatures` | Die Liste der Signaturen der einzelnen Dateien als Schlüssel-Wert-Paare, wobei der Schlüssel der Dateipfad ist und der Wert die Signatur. |
| `format`         | Die Kennung für das Format dieser Datei.                                                                                                  |
| `hostname`       | Der Name der Maschine, auf der die Signatur durchgeführt wurde.                                                                           | 
| `publicKey`      | Der öffentliche Schlüssel.                                                                                                                |
| `signatureType`  | Der Typ der Signatur.                                                                                                                     | 
| `timestamp`      | Der Zeitpunkt, zu dem die Signatur durchgeführt wurde.                                                                                    | 

### Formatkennung

Die Formatkennung gibt an, welches Format die Datei benutzt.
Zur Zeit ist nur ein Wert definiert: `1`.
Dieser Wert bedeutet, dass die Datei den Aufbau hat, der hier beschrieben ist.

### Signaturtyp

Der Signaturtyp kann zwei Werte haben:

| Signaturtyp | Bedeutung                                                                                                                                                                                                                                            |
|:-----------:|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|     `1`     | Die Signaturen sind mit dem Algorithmus [Ed25519](https://en.wikipedia.org/wiki/EdDSA#Ed25519) erstellt.                                                                                                                                             |
|     `2`     | Die Signaturen sind mit dem Algorithmus ECDSAP521 erstellt, also [ECDSA](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm) mit der Kurve [P-521](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-186.pdf). |

Beide Verfahren benutzen elliptische Kurven.
Weiteres ist in der Datei [Technische_Spezifikation.md](Technische_Spezifikation.md) zu finden.

### Zeitstempel

Der Zeitstempel liegt im [ISO 3339](https://datatracker.ietf.org/doc/html/rfc3339)-Format vor: `JJJJ-MM-TT hh:mm:ss +hh:mm`.

Der Zeitstempel beginnt mit dem Datum nach dem [gregorianischen Kalender](https://de.wikipedia.org/wiki/Gregorianischer_Kalender) in der Form Jahr (vierstellig), Monat (zweistellig), Tag (zweistellig).
Die einzelnen Teile sind durch ein Minuszeichen (`-`) voneinander getrennt.
Alle Teile sind mit führenden Nullen versehen.
Z.B. wird der Monat März als `03` dargestellt.

Danach folgt ein Leerzeichen zur Trennung.

Danach folgt die Uhrzeit in der Form Stunde (zweistellig, 24&nbsp;Stunden), Minute (zweistellig), Sekunde (zweistellig).
Die einzelnen Teile sind durch einen Doppelpunkt (`:`) voneinander getrennt.
Alle Teile sind mit führenden Nullen versehen.

Danach folgt ein Leerzeichen zur Trennung.

Danach folgt die Zeitzone in der Form Vorzeichen (`+`, `-`), Stunden (zweistellig), Doppelpunkt (`:`), Minuten (zweistellig).
Der Wert gibt an, um welche Zeitdifferenz sich die Zeitzone, auf die sich der angegebene Zeitstempel bezieht, von der [UTC](https://de.wikipedia.org/wiki/Koordinierte_Weltzeit) unterscheidet.
Hier ein paar Beispiele für den Wert der Zeitzonenangabe:

| Ort               | Abschnitt  | Zeitzone |
|-------------------|------------|----------|
| London            | Normalzeit | `+00:00` |
| London            | Sommerzeit | `+01:00` |
| Frankfurt am Main | Normalzeit | `+01:00` |
| Frankfurt am Main | Sommerzeit | `+02:00` |
| Los Angeles       | Normalzeit | `-08:00` |
| Los Angeles       | Sommerzeit | `-07:00` |
| Mumbai            | Normalzeit | `+05:30` |

### Kodierung von Binärwerten

Der öffentliche Schlüssel, die Dateisignaturen und die Signatur der Signaturendatei sind Binärwerte, die ähnlich dem [wort-sicheren Base32-Verfahren](https://en.wikipedia.org/wiki/Base32#Word-safe_alphabet) kodiert sind.

In dieser Base32-Kodierung steht ein Zeichen für 5 Bit aus dem Binärwert.
Das benutzte Alphabet lautet: `3479BCDFGHJLMRQSTVZbcdfghjmrstvz`.
Jedes Zeichen steht für einen Binärwert.
Die `3` steht für `0`, das `z`für `31` und die anderen Zeichen für die Werte dazwischen.

Sie enthält keine Vokale, so dass nicht versehentlich Worte erzeugt werden, die eine Bedeutung haben.
Gleichzeitig fehlen die Zeichen `6`, `X` und `x`, deren wiederholte Verwendung unerwünscht ist.
Ebenso enthält es keine Sonderzeichen, so dass man jeden so kodierten Wert in einem Editor mit einem Doppelklick vollständig markieren kann.

### Dateipfade

Die Dateipfade im Feld `fileSignatures` sind relative Dateipfade.
Sie beziehen sich also immer auf das aktuelle Verzeichnis.
Als Pfadtrenner wird immer der Schrägstrich (`/`) verwendet, unabhängig vom verwendeten Betriebssystem.
Die Dateipfade sind in UTF-8 kodiert.

## Regeln

Die JSON-Datei muss beim Einlesen auf formelle Fehler geprüft werden.
Dabei gelten folgende Regeln:

- Es **müssen** alle Felder vorhanden sein.
- Es **dürfen keine** zusätzlichen Felder vorhanden sein.

Sollte mindestens ein Feld fehlen oder mindestens ein zusätzliches Feld vorhanden sein, wird die Verarbeitung abgebrochen.

Nach dem Einlesen wird zuerst die Signatur `dataSignature` der Gesamtdatei geprüft.
Sollte die Gesamtsignatur nicht zu den restlichen Daten der Datei passen, wird die Verarbeitung abgebrochen.

## Beispiel

Die obigen Beschreibungen werden im folgenden Beispiel mit konkreten Werten erläutert.

Nehmen wir an, dass am 25.02.2024 um 13:37:22 auf dem Rechner "BuildHost", der sich in Mumbai befindet, für die folgenden Dateien eine Signaturendatei erstellt wurde:
`common.go`, `filesigner`, `filesigner.exe`, `maphelper/map_helper.go` und `set/set.go`.

Die Signaturendatei würde dann in etwa folgendermaßen aussehen:

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

Die Werte sind mit den obigen Erklärungen leicht nachzuvollziehen.

Als Beispiel ist hier der dekodierte Wert für den öffentlichen Schlüssel angegeben.
Der Base32-kodierte Wert `HxVJVrrjQcgfvhPxJ45chrQrRCFWmgJ5JH8JGMv6xxj23xjH8P52` entspricht den folgenden Byte-Werten: `5f e2 c8 f3 98 7d 2d 5e dd df 60 87 4b f1 fc 82 13 2c d9 83 62 cc c5 37 a4 ff f0 00 ff 0b 33 86`.
