package main

import "fmt"

// Definiere ein Struct Typ namens "Auto"
// - Der Typ hat folgende vier int Felder:
//   - batteriestand
//   - batterieverbrauch
//   - geschwindigkeit
//   - distanz
type Auto struct {
	batteriestand     int
	batterieverbrauch int
	geschwindigkeit   int
	distanz           int
}

// Definiere ein Struct Typ namens "Rennstrecke"
// - Der Typ hat ein int Feld namens "distanz"
type Rennstrecke struct {
	distanz int
}

func main() {
	fmt.Println(baueAuto(50, 10))
	fmt.Println(fahren(baueAuto(50, 10)))
	fmt.Println(prüfeFahrt(Auto{batteriestand: 0}, Rennstrecke{distanz: 10}))
	fmt.Println(prüfeFahrt(Auto{batteriestand: 100, batterieverbrauch: 1, geschwindigkeit: 0}, Rennstrecke{distanz: 10}))
	fmt.Println(prüfeFahrt(Auto{batteriestand: 100, batterieverbrauch: 1, geschwindigkeit: 10}, Rennstrecke{distanz: 10}))
}

// Definiere eine Funktion "baueAuto"
// - Die Funktion nimmt zwei Parameter
//   - Die Geschwindigkeit des Autos
//   - Den Batterieverbrauch des Autos
//   - Die Funktion gibt ein Auto zurück welches die Felder `batterieverbrauch` und `geschwindigkeit`
//     auf den Wert setzt der als Parameter übergeben wurde. `distanz` soll auf 0 und der `batteriestand` auf 100 sein
func baueAuto(geschwindigkeit int, batterieverbrauch int) Auto {
	return Auto{
		batterieverbrauch: batterieverbrauch,
		geschwindigkeit:   geschwindigkeit,
		distanz:           0,
		batteriestand:     100,
	}
}

// Definiere eine Funktion 'fahren'.
// - Die Funktion nimmt einen Parameter:
//   - Ein Auto
//
// - Die Funktion gibt ein Auto zurück mit veränderten Feldern:
//   - Der `batteriestand` wurde reduziert um den Wert in `batterieverbrauch`
//   - Die `distanz` wurde erhöht um den Wert in `geschwindigkeit`
func fahren(auto Auto) Auto {
	auto.batteriestand -= auto.batterieverbrauch
	auto.distanz += auto.geschwindigkeit
	return auto
}

// Definiere eine Funktion 'prüfeFahrt'
// - Die Funktion nimmt zwei Parameter:
//   - Das Auto
//   - Die Rennstrecke
//
// - Die Funktion gibt 'true' zurück wenn das Auto mit seinem Batteriestand die Rennstrecke komplett bis ins Ziel fahren kann
func prüfeFahrt(auto Auto, rennstrecke Rennstrecke) bool {
	if auto.batteriestand == 0 && rennstrecke.distanz != 0 {
		return false
	}

	return (auto.batteriestand/auto.batterieverbrauch)*auto.geschwindigkeit >= rennstrecke.distanz
}
