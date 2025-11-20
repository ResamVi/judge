# Zahl erraten

Das Programm soll sich eine Zahl zwischen 0 und 9 aussuchen.
Der Benutzer wird so lange nach einer Eingabe gefragt bis er die Zahl ausgegeben hat

Am Ende soll die Anzahl der Versuche ausgegeben werden auf folgende Weise:
```
Anzahl der Versuche: 6
```

## Ziel

1. `for`-Loops kennengelernt
2. `break` und `continue`

## Code

{{Code}}

## Erinnerung

```
for i := 1; i < 10; i++ {
    fmt.Println(i)
}
```

Inkrement / Dekrement

```
a := 10
a++ // gleich wie: a += 1
```

Zufallszahlen

```
import (
	"fmt"
	"math/rand"
)

func main() {
	n := rand.Intn(100)
	fmt.Println(n)
}
```