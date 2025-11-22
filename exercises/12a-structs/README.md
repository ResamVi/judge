# Structs

Schreibe drei Funktionen die mit structs arbeiten werden:

- Eine Funktion **kontoÖffnen** nimmt einen Namen und das Startguthaben und returned ein struct vom Typ `Konto` 
zurück mit den Feldern `Name` und `Guthaben` gesetzt auf die Werte die der Funktion übergeben worden

- Eine Funktion **werbungZeigen** welches bereits ein struct `Film` akzeptiert aber Felder noch nicht gesetzt sind

- Eine Funktion **flächeBerechnen** welches ein struct `Rechteck` akzeptiert und die Fläche des Rechtecks zurückgibt


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
