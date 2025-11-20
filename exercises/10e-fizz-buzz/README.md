# Fizz Buzz

Das Programm soll sich die Zahlen von 1 bis 20 Ausgeben.
- Wenn die Zahl durch drei teilbar ist soll es aber “Fizz” ausgeben
- Wenn die Zahl durch fünf teilbar ist soll es “Buzz” ausgeben
- Wenn die Zahl durch drei und fünf teilbar ist soll es “FizzBuzz” ausgeben

Erinnerung:
- i % 3 gibt den Rest einer ganzzahligen Divison (7 % 5 == 2)
- i % 3 == 0 bedeutet also … ?

## Ziel

1. `for`-Loops kennengelernt
2. Loops mit `if`s und `switch`s kombiniert

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