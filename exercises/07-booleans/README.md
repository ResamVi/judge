# Booleans

Schreibe folgende Funktionen:

- Erste Funktion: **eintrittErlaubt(hatTicket bool, istVIP bool) bool**  
Ausgabe: Der Eintritt ist erlaubt, wenn die Person ein Ticket hat oder VIP ist

- Zweite Funktion: **computerLäuft(eingesteckt bool, angeschalten bool) bool**  
Ausgabe: Der Computer läuft, wenn der Stecker eingesteckt ist und angeschalten wurde

- Dritte Funktion: **nameValide(enthältZahlen bool) bool**  
Ausgabe: Ein Name ist valide, wenn es keine Zahlen beinhaltet

- Vierte Funktion: **istHeiß(temperatur int) bool**  
Ausgabe: Die Temperatur ist heiß wenn sie über 100 Grad enthält

## Ziel

Folgende Konzepte kennengelernt
1. Boolean Logik

## Code

{{Code}}

## Erinnerung

Operationen
```
eins := true && true // true
zwei := true && false // false

drei := true || false // true
vier := false || false // false

fünf := !true // false
sechs := !false // true
```

Funktionen
```
func hi(name string) string {
    return "hi " + name
}
```
