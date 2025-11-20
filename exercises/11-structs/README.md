# Rennfahren

In dieser Aufgabe wird ein Rennen organisiert mit Batteriebetriebenen Autos.
Autos haben verschiedene Geschwindigkeiten und Akkulaufzeiten.

Wir schreiben ein Programm was die Autos auf verschiedenen Rennstrecken laufen lässt und prüft, ob sie ans Ende kommen.

## Neue Konzepte

1. Structs

## Code

{{Code}}

## Erinnerung

Neuen struct Typ definieren
```
type Person struct {
    name string
    alter int
}
```

Instanz des Typs erstellen
```
jules := Person {
    name: "Julien",
    alter: 28,
}
```

Felder lassen sich auslesen und anpassen
```
jules.alter = 29
fmt.Println("name: " + jules.name)
```
