# Funktionen

In einem Unternehmen berechnet sich das monatliche Gehalt eines Angestellten wie folgt:
Mindestlohn 400€ im Monat,
+ zuzüglich 20€ multipliziert mit der Anzahl der Beschäftigungsjahre,
+ zuzüglich 30€ für jedes Kind.
  Aufgabe: Schreibe eine Funktion 'berechneGehalt' mit erstem Parameter die Beschäftigungsjahre und zweitem Parameter die Anzahl der Kinder und welches das Gehalt zurückgibt

## Ziel

Folgende Konzepte kennengelernt
1. Funktionen

## Code

{{Code}}

## Erinnerung

Operationen
```
func Hello(name string) string {
    return "Hallo " + name
}

func HelloAndGoodbye(name string) (string, string) {
    return "Hallo " + name, "Tschüss " + name
}

func main() {
    greeting := Hello("Julien")
    fmt.Println(greeting)

    hello, goodbye := HelloAndGoodbye("Louisa")
    fmt.Println(hello)
    fmt.Println(goodbye)

    // Ausgabe:
    // Hallo Julien
    // Hallo Louisa
    // Tschüss Louisa
}
```
