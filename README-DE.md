# filesigner

Ein Programm, um eine einfache Möglichkeit zur Nutzung digitaler Signaturen zu haben.
Wie Zertifikate, nur besser und einfacher ;-).

[![Go Report Card](https://goreportcard.com/badge/github.com/xformerfhs/filesigner)](https://goreportcard.com/report/github.com/xformerfhs/filesigner)
[![License](https://img.shields.io/github/license/xformerfhs/filesigner)](https://github.com/xformerfhs/filesigner/blob/main/LICENSE)

## Einleitung

Die heutige IT-Welt besteht aus einem hochkomplexen Zusammenspiel vieler Komponenten:
Quellcode, Programme, Bibliotheken, Bilder, Dokumente, Konfigurationsdateien und noch vieles andere mehr werden erstellt, verteilt und in wichtigen Systemen eingesetzt.

Wie lässt sich sicherstellen, dass all diese Programme und Daten tatsächlich unverändert an ihren Zielorten ankommen?
Wie kann man erkennen, dass ein Programm, eine Datei, ein beliebiges Artefakt das ist, das durch den Build-Prozess erstellt wurde?

Das [NIST](https://www.nist.gov) fordert in seinem Secure Software Development Framework (SSDF) V1.1 ([NIST SP 800-218](https://nvlpubs.nist.gov/nistpubs/specialpublications/nist.sp.800-218.pdf)) im Punkt PS.2 die Bereitstellung eines Mechanismus zur Überprüfung der Software-Integrität.
Das [BSI](https://www.bsi.bund.de) fordert für den Grundschutz unter dem Punkt [APP.6.A4](https://www.bsi.bund.de/SharedDocs/Downloads/DE/BSI/Grundschutz/Kompendium_Einzel_PDFs_2021/06_APP_Anwendungen/APP_6_Allgemeine_Software_Edition_2021.pdf?__blob=publicationFile&v=1), das die Integrität von Installationsdateien geprüft werden muss.

Welche Methoden gibt es, um zu prüfen, ob die Artefakte die sind, die der Erzeuger erstellt hat?
Üblicherweise wird eines der beiden folgenden Verfahren benutzt:

- Hash-Werte der Artefakte 
- Digitale Signaturen

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

Die etablierten Methoden bieten entweder keine Sicherheit gegen einen Angreifer (Hash) oder sind sehr mühsam in der Anwendung (Signaturen mit Zertifikaten).

Daher wird hier ein Werkzeug angeboten, dass die Sicherheit von Signaturen bietet, ohne deren Komplexität.

## Beschreibung

Auch hier werden beliebige Dateien mit einer [digitalen Signatur](https://de.wikipedia.org/wiki/Digitale_Signatur) versehen.
Für eine [digitale Signatur](https://de.wikipedia.org/wiki/Digitale_Signatur) benötigt man ein [asymmetrisches Verschlüsselungsverfahren](https://de.wikipedia.org/wiki/Asymmetrisches_Kryptosystem) und eine [Hash-Verfahren](https://de.wikipedia.org/wiki/Hashfunktion).
Für die asymmetrische Verschlüsselung benötigt man ein Schlüsselpaar aus einem privaten und einem öffentlichen Schlüssel.
Über die zu signierenden Daten wird der Wert der Hash-Funktion berechnet und dieser dann mit dem privaten Schlüssel des Signaturerstellers mit dem asymmetrischen Verschlüsselungsverfahren verschlüsselt. 
Der öffentliche Schlüssel wird bekannt gegeben, damit man mit ihm die Signaturen prüfen kann.
Für die Prüfung wird die digitale Signatur mit dem öffentlichen Schlüssel entschlüsselt.
Dadurch erhält man den Hash-Wert des Signaturerstellers.
Diesen vergleicht man mit dem Hash-Wert, den man selbst ermittelt hat.
Stimmen beide überein, wurden die Daten nicht verändert.
Bei den herkömmlichen Verfahren werden alle Signaturen immer mit demselben Schlüsselpaar erzeugt, dessen privater Schlüssel sehr aufwändig vor unbefugtem Zugriff geschützt werden muss.

Das Neue an dem hier vorliegenden System ist, dass jeder Signaturblock mit einem neuen Schlüsselpaar durchgeführt wird.
Der private Schlüssel wird also nach jedem Signaturvorgang gelöscht und nirgends gespeichert.
Daher kann er auch nicht gestohlen und von Angreifern missbraucht werden.

Eine ausführliche Beschreibung ist in dem Dokument [Konzept.md](doc/de/Konzept.md) zu finden.

Der Nutzer der Artefakte benötigt für die Integritätsprüfung nur den öffentlichen Schlüssel des Erstellers.
In einem herkömmlichen System erhält er diesen durch das Zertifikat des Erstellers.
Man kann diesen öffentlichen Schlüssel jedoch auch auf andere Weise an den Nutzer weitergeben.
Es muss kein Zertifikat sein.
Es kann eine signierte E-Mail sein, eine vertrauenswürdige Web-Seite, ein Fax, ein Brief oder ein anderes weiteres Medium.
Das Wesentliche an dem hier vorgestellten System ist, dass der öffentliche Schlüssel und einige weitere Signaturinformationen auf einem anderen Weg, als die Signatur selbst bereitgestellt werden.

Die zu übermittelnden Informationen sind:

- Eine Kontext-Id
- Eine Id des Schlüssels
- Der Zeitstempel der Signaturerstellung
- Der Name des Rechners, auf dem die Signatur erstellt wurde

Für jede zu signierende Datei wird mit diesen Informationen und dem jeweiligen Dateiinhalt eine Signatur erstellt.

Schließlich wird noch die Signaturdatei selbst signiert.

Die Kontext-Id dient dazu, dem Nutzer eine Information mitzugeben, um was es sich bei den signierten Dateien handelt.
Bei Software könnte es sich zum Beispiel um die Versionsnummer handeln.
Es handelt sich also um so etwas, wie ein Thema.

Für die Signaturen selbst wird entweder das Verfahren [Ed25519](https://en.wikipedia.org/wiki/EdDSA#Ed25519) oder das Verfahren ECDSAP521 benutzt, also [ECDSA](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm) mit der Kurve [P-521](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-186.pdf). 

Die Signaturen werden in einem speziellen [Base32-Verfahren](https://en.wikipedia.org/wiki/Base32) kodiert.
Die Kodierung enthält keine Vokale, so dass nicht zufälligerweise echte Worte entstehen können.
Weiterhin gibt es keine leicht verwechselbaren Zeichen, wie `0` und `O` oder `1` und `l`.
Außerdem enthält die Kodierung keine Sonderzeichen, so dass ein kodierter Wert mit einem Doppelklick vollständig markiert werden kann.

Dies wird weiter unten an einem Beispiel dargestellt.

> [!IMPORTANT]
> Diese Software ist noch im Aufbau.
> Dies wird durch die Ziffer `0` an der ersten Stelle der Version angezeigt.
> Die Verfahren und Strukturen sind noch nicht endgültig festgelegt, auch wenn sie schon einige Prüfungen hinter sich haben.

## Aufrufe

### Signierung

Der Aufruf zur Signierung sieht folgendermaßen aus:

```
filesigner sign {contextId} [-a|--algorithm {algorithm}] [-i|--include-file {pattern}] [-x|--exclude-file {pattern}] [-I|--include-dir {pattern}] [-X|--exclude-dir {pattern}] [-f|--from-file {file}] [-m|--name {name}] [-r|--recurse] [-s|--stdin] [files...]
```

Die einzelnen Teile haben die folgenden Bedeutungen:

| Teil           | Bedeutung                                                                                                                                                                  |
|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `contextId`    | Ein beliebiger Text, der benutzt wird, um die Signatur von einem Thema abhängig zu machen.                                                                                 |
| `algorithm`    | Die Spezifikation der Signaturmethode. Entweder [`ed25519`](https://en.wikipedia.org/wiki/EdDSA) oder `ecdsap521`. Wird der Typ nicht angegeben, wird `ed25519` verwendet. |
| `exclude-dir`  | Spezifikation der Verzeichnisse, die nicht signiert werden sollen.                                                                                                         |
| `exclude-file` | Spezifikation der Dateien, die nicht signiert werden sollen.                                                                                                               |
| `from-file`    | Die zu bearbeitenden Dateinamen werden aus der angegebenen Datei gelesen, die einen Dateinamen pro Zeile enthalten muss.                                                   |
| `include-file` | Spezifikation der Dateien, die signiert werden sollen.                                                                                                                     |
| `include-dir`  | Spezifikation der Verzeichnisse, die signiert werden sollen.                                                                                                               |
| `name`         | Die Signaturendatei hat den Namen `{name}-signatures.json`. Die Voreinstellung für den Namen ist `filesigner`.                                                             |
| `recurse`      | Es werden auch Unterverzeichnisse bearbeitet.                                                                                                                              |
| `stdin`        | Die zu bearbeitenden Dateinamen werden von der Standardeingabe gelesen, die einen Dateinamen pro Zeile enthalten muss.                                                     |
| `files`        | Eine Liste von Dateinamen, die mit Leerzeichen getrennt sind.                                                                                                              |

Folgendes ist wichtig zu wissen:

* Alle exclude/include-Optionen durchlaufen das aktuelle Verzeichnis und alle Unterverzeichnisse, wenn `--recurse` angegeben ist.
* Alle exclude/include-Optionen müssen genau eine Dateispezifikation als Wert haben.
* In include/exclude-Optionen können Platzhalter (`*`, `?`) benutzt werden.
* Wenn sowohl Dateinamen als auch include-Optionen angegeben sind, werden sie zusammengefasst.
* Wenn sowohl Dateinamen als auch exclude-Optionen angegeben sind, werden Dateinamen, die zu einer exclude-Option passen, nicht signiert.
* Wenn in der Dateiliste Namen mit Wildcards enthalten sind, werden sie so behandelt, als ob sie in einer `--include-file`-Option angegeben wären.
* Eine include-Option schließt alle Objekte aus, die nicht in einer include-Option benannt werden.
* Unter Linux müssen Wildcards in einfache Anführungszeichen (`'`) oder doppelte Anführungszeichen (`"`) eingeschlossen werden oder mit einem vorangestellten \\ versehen werden (z.B.. `--exclude-dir .\*` um alle Verzeichnisse auszuschließen, die mit einem `.` beginnen).

> [!IMPORTANT]
> Die Signaturendatei wird **immer** ausgeschlossen und kann nicht signiert werden.
> Sie enthält bereits eine Signatur.

Der Aufruf erzeugt eine Datei[^1], die folgendes Format hat:

[^1]: Der Dateiinhalt ist nur ein Beispiel.
  Es ist nicht möglich, mit dieser Datei die Integrität von Dateien in diesem Repository zu verifizieren.

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

Es handelt sich um eine `json`-Datei, deren Aufbau und die Bedeutung der Felder in [Dateiformat.md](doc/de/Dateiformat.md) beschrieben ist.
Eine ausführliche Beschreibung der verschiedenen Berechnungen und Datenformate ist in [Technische_Spezifikation.md](doc/de/Technische_Spezifikation.md) zu finden.

Die Rückgabe-Codes können sein:

| Code | Bedeutung                    |
|------|------------------------------|
| `0`  | Verarbeitung erfolgreich     |
| `1`  | Fehler in der Befehlszeile   |
| `2`  | Warnung bei der Verarbeitung |
| `3`  | Fehler bei der Verarbeitung  |

### Verifizierung

Der Aufruf zur Verifizierung sieht folgendermaßen aus:

```
filesigner verify [-m|--name {name}]
```

Die einzelnen Teile haben die folgenden Bedeutungen:

| Teil           | Bedeutung                                                                                                                                                                  |
|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`         | Die Signaturendatei hat den Namen `{name}-signatures.json`. Die Voreinstellung für den Namen ist `filesigner`.                                                             |

> [!IMPORTANT]
> Weitere Parameter sind nicht erlaubt und führen zu einer Fehlermeldung.

Das Programm liest die Signaturendatei ein und prüft, ob die dort genannten Dateien vorhanden sind und ob deren Signaturen zu den aktuellen Inhalten passen.

Die Rückgabe-Codes sind dieselben, wie bei der Signierung.

## Programme

| OS      | Programm         |
|---------|------------------|
| Windows | `filesigner.exe` |
| Linux   | `filesigner`     |

Das Linux-Programm ist auf jedem beliebigen Linux ausführbar.

## Beispiel

### Signierung

Angenommen, für die Artefakte `filesigner`, `filesigner.exe`, alle `Go`-Dateien und alle Dateien, die mit dem Wort `go`beginnen, sollen für die Version `1.7.11` einer Anwendung Signaturen erzeugt und überprüft werden.

Dann erzeugt man die Signatur mit folgendem Aufruf:

```
filesigner sign project1711 -if *.go -if filesign*
```

Das Programm erzeugt dann die folgende Ausgabe auf der Konsole:

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

Der Rückgabe-Code ist 0.

### Verifizierung

Zur Verifizierung der Signaturen benötigt man einen vertrauenswürdigen Ort, an dem die Id des öffentlichen Schlüssels, der Zeitstempel der Signatur und der Name des Signatur-Rechners veröffentlicht werden.
Dabei kann es sich um eine signierte E-Mail, eine Website, eine Datenbank oder einen anderen sicheren Ort handeln, der als vertrauenswürdig eingestuft wird.

Zur Verifizierung ruft man das Programm folgendermaßen auf:

```
filesigner verify project1711
```

Das Programm erzeugt dann die folgende Ausgabe auf der Konsole:

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

Der Rückgabe-Code ist 0.

Die verifizierende Person prüft, ob die angezeigte Id des öffentlichen Schlüssels, der Zeitstempel der Signatur und der Name des Signatur-Rechners der Signatur mit denen übereinstimmen, die am vertrauenswürdigen Ort gespeichert sind.
Ist dies nicht der Fall, wird die Signatur als ungültig angesehen und die Dateien dürfen nicht als vertrauenswürdig angesehen werden!

Sollte, als weiteres Beispiel, die Datei `filesigner` manipuliert worden sein, würde folgende Ausgabe erscheinen:

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

Der Rückgabe-Code ist 3.

Sollte z.B. die Signaturdatei manipuliert worden sein, würde folgende Ausgabe erscheinen:

```
2024-03-05 15:50:04 +01:00  15  I  filesigner V0.83.1 (go1.24.3, 8 cpus)
2024-03-05 15:50:04 +01:00  51  I  Reading signatures file 'filesigner-signatures.json'
2024-03-05 15:50:04 +01:00  54  E  Signatures file has been modified
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

Frank Schwab ([Mail](mailto:github.sfdhi@slmails.com "Mail"))

## Lizenz

Das Programm ist unter der [Apache-Lizenz V2](https://www.apache.org/licenses/LICENSE-2.0.txt) veröffentlicht.
