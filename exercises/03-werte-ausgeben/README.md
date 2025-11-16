# Werte ausgeben

Werte repräsentieren Daten. Wir haben 
- Ganze Zahlen (`42`, `0`, `-1`)
- Kommazahlen (`4.20`, `3.131313`)
- Zeichenketten genannt: 'Strings' (`"Hallo"`, `"Eins und zwei"`)
- Wahrheitswerte genannt: 'Booleans' (`true`, `false`)

## Ziel

Folgende Konzepte kennengelernt
1. Ganze Zahlen, Kommazahlen, Wörter und Wahrheitswerte

## Code

{{Code}}

## Ratschläge

- Kommazahlen werden mit Punkten geschrieben `3.14159`
- fmt.Println() akzeptiert mehrere Werte, die man mit Komma trennen muss (überhaupt nicht verwirrend!)
```go
fmt.Println("abc", "123") // Ausgabe: abc 123
```