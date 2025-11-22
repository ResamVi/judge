# Errors
Schreibe ein Programm, welches nach einem Passwort fragt und dann prüft, ob es mehr als 8 Zeichen hat.
Falls nicht, soll die Funktion einen Fehler zurückgeben der sagt "Passwort ist zu kurz"

## Neue Konzepte

1. Errors

## Code

{{Code}}

## Erinnerung

Eine Funktion die einen error zurückgibt
```
func division(a float64, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}

	result := a / b
	return result, nil
}
```

Nach einem Error prüfen
```
	result, err := division(3, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
```
