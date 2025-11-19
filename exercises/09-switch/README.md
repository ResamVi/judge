# Switch

In dieser Übung simulieren wir den ersten Zug einer Blackjack-Runde.

Du erhältst zwei Karten und siehst welche Karte der Dealer offengelegt hat.
Alle Karten sind repräsentiert durch Strings wie "Ass", "König", "Drei", "Zwei", etc.

Die Kartenwerte in Blackjack sind:

```
Ass 	11 	Acht 	8
Zwei 	2 	Neun 	9
Drei 	3 	Zehn 	10
Vier 	4 	Bube 	10
Fünf 	5 	Dame 	10
Sechs 	6 	König 	10
Sieben 	7
```

(In Wirklichkeit zählen Asse entweder 11 oder 1 sind aber für die Einfachheit nur 11 in dieser Aufgabe)

## Ziel

1. `Switch`-Bedingungen kennengelernt

## Code

{{Code}}

## Erinnerung

Switch-Bedingungen
```
switch betriebssystem {
case "windows":
    fmt.Println("Programm ist eine .exe")
case "linux":
    fmt.Println("Programm ist eine ./binary")
case "macos":
    fmt.Println("Programm ist eine ./binary")
default:
    fmt.Println("Programm ist unbekannt")
}
```

Switch-Bedingung ohne Wert
```
age := 21

switch {
case age > 20 && age < 30:
    return "Du bist in deinen Zwanziger"
case age == 10:
    return "Dein Alter ist Zweistellig"
default:
    return "Cooles Alter!"
}
```
