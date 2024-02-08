# filesigner

Ein Programm, um eine einfache Möglichkeit zur Nutzung digitaler Signaturen zu haben.
Wie Zertifikate, nur besser und einfacher ;-).

## Einleitung

Wenn man Artefakte benutzt, benötigt man eine Methode, um zu prüfen, ob diese Artefakte die sind, die der Erzeuger erstellt hat.
Zurzeit gibt es keine Möglichkeit, zu überprüfen, ob an den abgelegten Artefakten etwas verändert wurde.
Es gibt verschiedene Möglichkeiten, eine solche Überprüfung durchzuführen:

- Nutzung von Hash-Werten
- Nutzung von digitalen Signaturen

Hash-Werte sind zwar leicht auszurechnen und zu prüfen, aber sie bieten nur Schutz gegen irrtümliche Änderungen, nicht gegen Angriffe.
Ein Angreifer, der die Artefakte ändern kann, kann auch die publizierten Hash-Werte ändern.

Digitale Signaturen bieten einen Schutz vor solchen Angriffen, da der Angreifer für eine Fälschung Zugriff auf den privaten Schlüssel der Signatur haben müsste.
Sie arbeiten üblicherweise mit Zertifikaten, doch diese sind schwierig handhabbar:

- Der private Schlüssel eines Zertifikats muss unbedingt stark geschützt werden.
Wird er nicht autorisierten Personen bekannt, können diese Signaturen für manipulierte Artefakte selbst erstellen.
- Es muss immer überprüft werden, ob das Zertifikat gültig ist und ob es inzwischen zurückgezogen wurde.
Das ist ausgesprochen mühsam und fehleranfällig.
- Zertifikate laufen ab und müssen regelmäßig erneuert werden.
Das bedeutet einen erheblichen organisatorischen Aufwand.

## ACHTUNG

Die jetzige Version trägt vorne noch eine `0` in der Versionsnummer.
Die Schnittstelle kann sich also noch ändern.
Jede konstruktive Rückmeldung zu diesem Programm ist willkommen.

## Beschreibung

Die hier vorliegende Lösung bietet eine digitale Signatur, ohne die Schwierigkeiten, die mit Zertifikaten verbunden sind.
Artefakte werden signiert.
Dafür wird ein Schlüsselpaar aus einem privaten und öffentlichen Schlüssel erzeugt.
Die Signaturen werden mit dem privaten Schlüssel ausgeführt.
Der öffentliche Schlüssel wird ausgegeben, damit man mit ihm die Signaturen prüfen kann.
Nach dem Signierungsprozess wird der private Schlüssel gelöscht.
Er wird nicht gespeichert und kann daher auch nicht gestohlen und von Angreifern missbraucht werden.
Eine Überprüfung der Signatur ist jedoch weiterhin durch den öffentlichen Schlüssel möglich.

Wie kann man sich nun dagegen schützen, dass ein Angreifer sein gefälschtes Artefakt ablegt und die Signatur mit dem passenden Programm erzeugt?

Bei der Veröffentlichung von Artefakten werden diese und die Signaturdatei abgelegt.
Der verwendete öffentliche Schlüssel wird den Abnehmern der Artefakte auf einem anderen Weg bekannt gemacht.
Diese können dadurch immer prüfen, ob die Signaturdatei auch die ist, die vom Veröffentlichungsteam herausgegeben wurde.

Dies wird weiter unten an einem Beispiel dargestellt.

## Aufrufe

### Signierung

Der Aufruf zur Signierung sieht folgendermaßen aus:

```
filesigner sign {contextId} [-a|--algorithm {algorithm}] [-i|--include-file {pattern}] [-x|--exclude-file {pattern}] [-j|--include-dir {pattern}] [-y|--exclude-dir {pattern}] [-f|--from-file {file}] [-s|--signatures-file {file}] [-r|--recurse] [-n|--stdin] [files...]
```

Die einzelnen Teile haben die folgenden Bedeutungen:

| Teil | Bedeutung |
|-------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `contextId`       | An arbitrary text used to make the signature depend on a topic, also called a "domain separator".                                                               |
| `algorithm`       | Specification of the signature method. Either [`ed25519`](https://en.wikipedia.org/wiki/EdDSA) or `ecdsap521`. If the type is not specified, `ed25519` is used. |
| `exclude-dir`     | Specification of directories to exclude.                                                                                                                        |
| `exclude-file`    | Specification of files to exclude.                                                                                                                              |
| `from-file`       | Read file names to process from the specified file. There is one file name per line.                                                                            |
| `include-dir`     | Specification of directories to include.                                                                                                                        |
| `include-file`    | Specification of files to include.                                                                                                                              |
| `recurse`         | Descend also into subdirectories.                                                                                                                               |
| `signatures-file` | Put signatures in the specified file. Default is `signatures.json`.                                                                                             |
| `stdin`           | Read file names to process from the standard input. There is one file name per line.                                                                            |
| `files`           | A blank-separated list of files to sign.                                                                                                                        |

Folgendes ist wichtig zu wissen:

* The exclude/include options scan the current directory
  and the subdirectories if `--recurse` is specified.
* All exclude/include options take one specification.
* Wildcards (`*`, `?`) may be used in include/exclude options.
* An include option excludes all objects that are not included.
* If both, files and includes are specified, they are combined.
* If both, files and excludes are specified, files that match an exclude specification are not processed.
* If wildcards are specified in the files list, they are treated as if they are values in `--include-file` options.
* On Linux, wildcards need to be put in quotes (`'`) or double quotes (`"`) or escaped by a \\ (like e.g. `--exclude-dir .\*` to exclude all directories starting with `.`).

```
filesigner sign {contextId} [-type {type}] [-if|-include-file {mask}] [-xf|-exclude-file {mask}] [-id|-include-dir {mask}] [-xd|-exclude-dir {mask}] [-no-subdirs]
```

Die einzelnen Teile haben die folgenden Bedeutungen:

| Teil              | Bedeutung                                                                                                                                                                  |
|-------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `contextId`       | Ein beliebiger Text, der benutzt wird, um die Signatur von einem Thema abhängig zu machen.                                                                                 |
| `algorithm`       | Die Spezifikation der Signaturmethode. Entweder [`ed25519`](https://en.wikipedia.org/wiki/EdDSA) oder `ecdsap521`. Wird der Typ nicht angegeben, wird `ed25519` verwendet. |
| `exclude-dir`     | Spezifikation der Verezichnisse, die nicht signiert werden sollen.                                                                                                         |
| `exclude-file`    | Spezifikation der Dateien, die nicht signiert werden sollen.                                                                                                               |
| `from-file`       | Die zu bearbeitenden Dateinamen werden aus der angegebenen Datei gelesen, die einen Dateinamen pro Zeile enthalten muss.                                                   |
| `include-file`    | Spezifikation der Dateien, die signiert werden sollen.                                                                                                                     |
| `include-dir`     | Spezifikation der Verzeichnisse, die signiert werden sollen.                                                                                                               |
| `recurse`         | Es werden auch Unterverzeichnisse bearbeitet.                                                                                                                              |
| `signatures-file` | Die Signaturen werden in die angegebenene Datei geschrieben. Die Voreinstellung ist `signatures.json`.                                                                     |
| `stdin`           | Die zu bearbeitenden Dateinamen werden von der Standardeingabe gelesen, die einen Dateinamen pro Zeile enthalten muss.                                                     |
| `files`           | Eine Liste von Dateinamen, die mit Leerzeichen getrennt sind.                                                                                                              |

Folgendes ist wichtig zu wissen:

* Alle exclude/include-Optionen durchlaufen das aktuelle Verzeichnis und 
alle Unterverzeichnisse, wenn `--recurse` angegeben ist.
* Alle exclude/include-Optionen müssen genau eine Dateispezifikation als Wert haben.
* In include/exclude-Optionen können Platzhalter (`*`, `?`) benutzt werden.
* Wenn sowohl Dateinamen, als auch include-Optionen angegeben sind, werden sie zusammengefasst.
* Wenn sowohl Dateinamen, als auch exclude-Optionen angegeben sind, werden Dateinamen, die zu einer
exclude-Option passen, nicht signiert.
* Wenn in der Dateiliste Namen mit Wildcards enthalten sind, werden sie so behandelt, als ob sie in einer `--include-file`-Option angegeben wären.
* Eine include-Option schließt alle Objekte aus, die nicht in einer include-Option benannt werden.
* Unter Linux müssen Wildcards in einfache Anführungszeichen (`'`) oder doppelte Anführungszeichen (`"`) eingeschlossen werden oder mit einem vorangestellten \\ versehen werden (z.B.. `--exclude-dir .\*` um alle Verzeichnisse auszuschließen, die mit einem `.` beginnen).

Die Signaturendatei wird **immer** ausgeschlossen und kann nicht signiert werden.
Sie enthält bereits eine Signatur.

Der Aufruf erzeugt eine Datei mit dem Namen `signatures.json`, die folgendes Format hat:

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

Es handelt sich um eine `json`-Datei mit den Feldern

| Feld             | Bedeutung                                                                                                                                                                                                                                                                                                          |
|------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Format`         | Eine Zahl mit der Kennung des Formats.                                                                                                                                                                                                                                                                             |
| `PublicKey`      | Der verwendete öffentliche Schlüssel in einer speziellen Base32-Kodierung.                                                                                                                                                                                                                                         |
| `Timestamp`      | Der Zeitpunkt, zu dem die Signatur durchgeführt wurde.                                                                                                                                                                                                                                                             | 
| `Hostname`       | Der Name der Maschine, auf der die Signatur durchgeführt wurde.                                                                                                                                                                                                                                                    | 
| `SignatureType`  | Typ der Signatur<br/>          1: [Ed25519](https://en.wikipedia.org/wiki/EdDSA) mit [SHA3-512](https://en.wikipedia.org/wiki/SHA-3)-Hash<br/>2: [EcDsa](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm) mit der Kurve [P521](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-186.pdf) und [SHA3-512](https://en.wikipedia.org/wiki/SHA-3)-Hash | 
| `FileSignatures` | Die Signaturen der einzelnen Dateien als Schlüssel-Wert-Paare, wobei der Schlüssel der Dateiname ist und der Wert die Signatur in derselben speziellen Base32-Kodierung, wie beim PublicKey.                                                                                                                       |
| `DataSignature`  | Die Signatur über die einzelnen Teile dieser Datei in der speziellen Base32-Kodierung.                                                                                                                                                                                                                             |

In die Signaturen fließen sowohl der Inhalt der Dateien, als auch der Signaturzeitpunkt und der Maschinenname ein.

Die Rückgabe-Codes können sein:

| Code | Bedeutung                     |
|------|------------------------------|
| `0`  | Verarbeitung erfolgreich     |
| `1`  | Fehler in der Befehlszeile   |
| `2`  | Warnung bei der Verarbeitung |
| `3`  | Fehler bei der Verarbeitung  |

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

Binärwerte sind in einer speziellen wort-sicheren [Base32-Kodierung](https://en.wikipedia.org/wiki/Base32) abgelegt.
Die Besonderheit besteht darin, dass das verwendete Alphabet keine Vokale, keine leicht zu verwechselnden Zeichen wie `0` und `O` oder `1` und `l` und keine Sonderzeichen enthält.
Dadurch können die so kodierten Werte mit einem Doppelklick markiert werden, es können nicht versehentlich echte Worte entstehen und beim Vorlesen kann man keine Zeichen verwechseln.

## Beispiel

### Signierung

Angenommen, für die Artefakte `filesigner`, `filesigner.exe`, alle `Go`-Dateien und alle Dateien, die mit dem Wort `go`beginnen, sollen für die Version `1.7.11` einer Anwendung Signaturen erzeugt und überprüft werden.

Dann erzeugt man die Signatur mit folgendem Aufruf:

```
filesigner sign projekt1711 -if *.go -if filesign*
```

Das Programm erzeugt dann die folgende Ausgabe auf der Konsole:

```
2023-06-13 12:51:25 +02:00  17  I  filesigner V0.50.0 (go1.21.6, 12 cpus)
2023-06-13 12:51:25 +02:00  37  I  Context id         : projekt1711
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

Der Rückgabe-Code ist 0.

### Verifizierung

Zur Verifizierung der Signaturen benötigt man einen vertrauenswürdigen Ort, an dem die Id des öffentlichen Schlüssels, der Zeitstempel der Signatur und der Name des Signatur-Rechners veröffentlicht werden.
Dabei kann es sich um eine signierte E-Mail, eine Website, eine Datenbank oder einen anderen sicheren Ort handeln, der als vertrauenswürdig eingestuft wird.

Zur Verifizierung ruft man das Programm folgendermaßen auf:

```
filesigner verify projekt1711
```

Das Programm erzeugt dann die folgende Ausgabe auf der Konsole:

```
2023-06-13 12:51:52 +02:00  17  I  filesigner V0.50.0 (go1.21.6, 12 cpus)
2023-06-13 12:51:53 +02:00  55  I  Context id         : projekt1711
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

Der Rückgabe-Code ist 0.

Die verifizierende Person prüft, ob die angezeigte Id des öffentlichen Schlüssels, der Zeitstempel der Signatur und der Name des Signatur-Rechners der Signatur mit denen übereinstimmen, die am vertrauenswürdigen Ort gespeichert sind.
Ist dies nicht der Fall, wird die Signatur als ungültig angesehen und die Dateien dürfen nicht als vertrauenswürdig angesehen werden!

Sollte, als weiteres Beispiel, die Datei `filesigner` manipuliert worden sein, würde folgende Ausgabe erscheinen:

```
2023-06-13 12:51:52 +02:00  17  I  filesigner V0.50.0 (go1.21.6, 12 cpus)
2023-06-13 12:51:53 +02:00  55  I  Context id         : projekt1711
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

Der Rückgabe-Code ist 3.

Sollte z.B. die Signaturdatei manipuliert worden sein oder die falsche Kontext-Id angegeben worden sein, würde folgende Ausgabe erscheinen:

```
2023-06-13 12:51:52 +02:00  17  I  filesigner V0.50.0 (go1.21.6, 12 cpus)
2023-06-13 12:51:53 +02:00  55  I  Context id         : projekt1711
2023-06-13 12:51:53 +02:00  56  I  Public key id      : 85R3-VZPX-JRV8-RN6R-G0L1-SV4U-NW
2023-06-13 12:51:53 +02:00  57  I  Signature timestamp: 2023-06-13 12:51:25 +02:00
2023-06-13 12:51:53 +02:00  58  I  Signature host name: MDXN01022044
2023-06-13 12:54:56 +02:00  53  E  Signature file has been tampered with or wrong context id
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

## Kontakt

Frank Schwab ([Mail](mailto:frank.schwab@live.de "Mail"))

## Lizenz

Das Programm ist unter der [Apache-Lizenz V2](https://www.apache.org/licenses/LICENSE-2.0.txt) veröffentlicht.
