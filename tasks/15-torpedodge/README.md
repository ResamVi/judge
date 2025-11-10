# Kartentricks

Wir arbeiten mit einem Stapel an Spielkarten und wollen diese in einem Programm etwas manipulieren.
Um die Sache einfacher zu machen gehen die Karten nur von 1 bis 10

## Neue Konzepte

1. Slices

## Code

{{Code}}

2. Downloade den Code mit:   

```
judge download 10
```


3. Uploade den Code mit:
```
judge upload 10
```

## Erinnerung

Wie man slices definiert
```
var leer []int                 
mitInhalt := []int{0,1,2,3,4,5}
```

Einzelne Felder eines Slices füllen

```
mitInhalt[1] = 5
x := mitInhalt[1] // x ist 5
```

Neue Elemente an ein Slice hängen mit `append`
```
a := []int{1, 3}
a = append(a, 4)
// => []int{1,3,4}
```

Länge eines Slices herausfinden
```
a := []int{1, 3, 5, 7}
laenge := len(a)
// => laenge ist 4
```
