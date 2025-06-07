# Blackjack

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

In Wirklichkeit zählen Asse entweder 11 oder 1 sind aber für die Einfachheit nur 11 in dieser Aufgabe

## Ziel

1. `Switch`-Bedingungen kennengelernt

## Code

{{Code}}

2. Downloade den Code mit:   

```
judge download 07
```


3. Uploade den Code mit:
```
judge upload 07
```

## Erinnerung

Die logischen Operatoren
```
Gleichheit              ==
Ungleichheit 	        !=
Weniger als 	        <
Gleich oder weniger     <=
Mehr als                >
Gleich oder mehr     	>=
```

Ergebnisse von logischen Operatoren sind Booleans
```
a != 4 // true
a > 5  // false

"apple" < "banana"  // true
"apple" > "banana"  // false
```

If-Bedingungen
```
if zahl > 0 {
    return "Positiv"
} else if zahl < 0 {
    return "Negativ"
} else {
    return "Null"
}
```
