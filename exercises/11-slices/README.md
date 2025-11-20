# Slices

Das Programm arbeitet mit einem Stapel von Spielkarten.

Zuerst:
- Schreib eine Funktion Namens **stapel** welche ein Stapel von 14 Karten zurückgibt
Beispielausgabe: `["1" "2" "3" "4" "5" "6" "7" "8" "9" "10" "J" "Q" "K" "A"]`

- Schreibe eine Funktion **karteNehmen**, welche einen Stapel nimmt und eine Zahl n für die n-te Karte und diese ausgibt

- Schreibe eine Funktion **karteAustauschen** welche einen Stapel, eine Zahl n und eine Karte nimmt und die Karte an der Stelle n ändert

- Schreibe eine Funktion **karteHinzufügen** welche einen Stapel und eine Karte nimmt und die Karte auf den Stapel legt

## Neue Konzepte

1. Slices als Container für mehrere Werte

## Code

{{Code}}

## Erinnerung

Deklaration
```
var meinSlice []int                       // Leer
meinSlice := []int{9, 10, 11, 12, 13, 14} // mit Daten gefüllt
}
```

Zugriff auf Elemente des Slices
```
fmt.Println(meinSlice[3]) // 12
fmt.Println(meinSlice[0]) // 9
fmt.Println(meinSlice[6]) // panic: runtime error: index out of range [6] with length 6
```

Element ändern in Slice
```
meinSlice[1] = 5
```

Element hinzufügen
```
a := []int{1, 3}
a = append(a, 4)
// => []int{1,3,4}
```