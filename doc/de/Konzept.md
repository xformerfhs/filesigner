# Konzept

## Beschreibung

In diesem Dokument wird das Konzept beschrieben.

## Motivation

Es gab und gibt immer wieder Angriffe auf Firmen und Privatpersonen, bei denen manipulierte Programme oder Dateien eingesetzt wurden.
In den letzten Jahren hat sich die Erkenntnis durchgesetzt, dass solche Manipulationen erkannt werden müssen, um die Sicherheit aller zu erhöhen.

Die [Open Source Security Foundation](https://openssf.org/) hat viele Projekte gestartet, um die Sicherheit bei Open-Source-Projekten zu erhöhen.
Eines davon ist das [Scorecard](https://github.com/ossf/scorecard/tree/main)-Projekt, das fordert, dass Artefakte digital signiert sind.

Im [SigStore](https://www.sigstore.dev/)-Projekt[^1] wird eine Methode und eine Infrastruktur eingerichtet, mit der man jederzeit prüfen kann, ob ein Container-Image tatsächlich das ist, das der Ersteller veröffentlicht hat.

[^1]: SigStore benutzt ähnliche Methoden. 
  Das hier vorliegende Verfahren wurde von 2021 bis 2024 ohne Kenntnis des SigStore-Projektes entwickelt.

Das [NIST](https://www.nist.gov) hat im Standard [SP 800-218](https://csrc.nist.gov/pubs/sp/800/218/final) das "Secure Software Development Framework (SSDF)" veröffentlicht, wo in der Maßnahme "PS 2.1" gefordert wird, dass der Ersteller einer Software Informationen zur Verfügung stellen muss, die es dem Benutzer erlauben, die Integrität der Software-Artefakte zu verifizieren.

Aber nicht nur bei Programm-Artefakten wie Programmen, Bibliotheken oder Konfigurationsdateien ist es wichtig, zu überprüfen, ob die Daten unverändert sind.
Auch bei vielen anderen Daten ist das nötig, wie Bildern, Dokumenten und jeder beliebigen Datenart.

## Hintergrund

Eine Integritätssicherung wird durch das Zusammenspiel zweier Komponenten erzeugt:

1. Über die Daten wird eine Prüfsumme ([Hash](https://de.wikipedia.org/wiki/Kryptographische_Hashfunktion)) berechnet.
2. Die Prüfsumme wird mit dem privaten Schlüssel eines [asymmetrischen Verschlüsselungsverfahrens](https://de.wikipedia.org/wiki/Asymmetrisches_Kryptosystem) verschlüsselt.

Diesen Wert nennt man eine [digitale Signatur](https://de.wikipedia.org/wiki/Digitale_Signatur).

Um die Signatur zu überprüfen, muss man die drei folgenden Schritte durchführen:

1. Berechnung der Prüfsumme (Hash) der lokalen Kopie der Daten.
2. Entschlüsselung der digitalen Signatur mit dem zum privaten Schlüssel gehörenden öffentlichen Schlüssel.
3. Vergleich der lokal berechneten Prüfsumme mit der entschlüsselten Prüfsumme.

Wenn beide Prüfsummen gleich sind, ist die digitale Signatur gültig.

Dieses Verfahren stellt sicher, dass

1. Die Daten nicht verändert wurden.
2. Der Ersteller der digitalen Signatur im Besitz des zum öffentlichen Schlüssel gehörenden privaten Schlüssels war.

Um nun sicherzustellen, dass der öffentliche Schlüssel, mit dem man die digitale Signatur überprüft, tatsächlich derjenige ist, der er behauptet zu sein, gibt es [Zertifikate](https://de.wikipedia.org/wiki/Digitales_Zertifikat).
In diesem attestiert eine vertrauenswürdige Stelle, dass der öffentliche Schlüssel zu einer bestimmten Entität gehört.
Das kann eine E-Mail-Adresse oder eine Web-Seite sein oder eine Person oder Firma.

## Schwierigkeiten mit dem Zertifikatssystem

Zertifikate haben leider sehr viele Probleme:

- Der private Schlüssel des Zertifikats muss permanent gespeichert und **sehr gut** gesichert sein.
  - Jeder, der an den privaten Schlüssel gelangt, kann im Namen des Eigentümers des Zertifikates Artefakte signieren.
  - Es ist sehr schwierig den privaten Schlüssel gut abzusichern und meist sehr teuer, da sie den Einsatz spezieller Hardware ([HSM](https://de.wikipedia.org/wiki/Hardware-Sicherheitsmodul)) nötig macht. 
- Fälschliche oder irrtümliche Ausstellung.
- Sie können zurückgezogen werden.
  - Es ist mühsam und fehleranfällig auf zurückgezogene Zertifikate zu prüfen.
- Sie haben ein Ende-Datum, nach dem sie nicht mehr gültig sind.

Besonders der letzte Punkt führt zu enormen Problemen:

- Abgelaufene Zertifikate führen zu Totalausfällen von IT-Verfahren.
- Wenn man eine signierte E-Mail erhält, ist das Zertifikat, das zu der Signatur gehört mit einer begrenzten Gültigkeit versehen.
Fast alle E-Mail-Programmierer machen den Fehler die Signatur einer E-Mail als "ungültig" anzuzeigen, wenn das dazugehörige Zertifikat abgelaufen ist.
Tatsächlich ist die Signatur gültig, denn das Zertifikat war zum Zeitpunkt der digitalen Signatur gültig.

Programmierer und Betriebsführer sind mit der korrekten Handhabung von Zertifikaten oft überfordert.

## Lösungsansatz

Das hier vorgestellte Signatursystem vermeidet die beiden größten Schwachpunkte von Zertifikate:

1. Ablaufdatum
2. Permanente Speicherung des privaten Schlüssels

Zum einen wurde die vorliegende Lösung inspiriert von der Verschlüsselungs-Software [age](https://github.com/FiloSottile/age) des Verschlüsselungsexperten Filippo Valsorda.
In age gibt es keine Zertifikate, die ablaufen.
Der öffentliche Schlüssel wird veröffentlicht und ist so lange gültig, wie er veröffentlicht ist.

Damit ist der erste Schwachpunkt behoben.

Die zweite Idee stammt von meinem hochgeschätzten Kollegen Florian Schäfer von der DB Systel GmbH:
Der private Schlüssel wird nach der Signierung **nicht** gespeichert.
Damit kann man ihn nicht mehr stehlen und nicht mehr missbrauchen.
Nur der öffentliche Schlüssel eines konkreten Signierungsvorganges wird an einem vertrauenswürdigen Ort veröffentlicht.
Mit ihm ist die Prüfung einer Signatur möglich.

Dasselbe Verfahren benutzt das SigStore-Projekt.
Auch [dort](https://docs.sigstore.dev/signing/overview/#verifying-identity-and-signing-the-artifact) wird ein kurzlebiges Schlüsselpaar erzeugt und der private Schlüssel direkt nach der Signierung gelöscht.
Nur der öffentliche Schlüssel wird beibehalten, um eine Überprüfung der Signatur zu ermöglichen.
Bei SigStore wird das Ergebnis in einem "Transparency Log" gespeichert, was nichts anderes ist, als ein vertrauenswürdiger Speicherort für das Signaturergebnis.

Das hier vorliegende Verfahren gibt keinen Speicherort vor.

Für jeden Signierungsvorgang werden die folgenden Informationen veröffentlicht:

1. Der öffentliche Schlüssel
2. Eine Kontext-Id
3. Der Zeitstempel des Signierungsvorganges
4. Der Rechnername, auf dem der Signaturvorgang durchgeführt wurde

Diese vier Informationen müssen an einem vertrauenswürdigen Ort veröffentlicht werden.

Die Kontext-Id gibt dabei genauere Informationen, um was es sich handelt.
Das kann eine Build-Nummer sein, eine Versionsinformation, eine Kennung für die Ziel- oder Quellumgebung.
Es ist ein beliebiges Freitext-Feld.

Alle diese Informationen werden zur Erstellung der Signatur benutzt.

Die Signaturen für die zu signierenden Dateien werden in einer Datei gesammelt.
Diese Datei ist selbst signiert.

Wenn nun jemand die Artefakte bekommt, kann er mit der Signaturendatei die Signaturen aller Dateien prüfen.

Im vorliegenden Repository wurde ein solches Verfahren in [Go](https://go.dev/) implementiert.

Es lässt sich in jeder beliebigen anderen Programmiersprache ebenfalls implementieren.
Das Format der Signaturendatei ist in [Dateiformat.md](Dateiformat.md) beschrieben.
Die technische Spezifikation der dort enthaltenen Werte ist in [Technische_Spezifikation.md](Technische_Spezifikation.md) beschrieben.
