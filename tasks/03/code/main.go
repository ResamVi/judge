package main

import "fmt"

var (
	ritterWach = true
	schützeWach = false
	prinzWach = true
)

func main() {
	fmt.Println(schnellangriff(ritterWach))
	fmt.Println(spionieren(ritterWach, schützeWach, prinzWach))
	fmt.Println(brieftaube(schützeWach, prinzWach))
	fmt.Println(retten(ritterWach, schützeWach, prinzWach))
}

// TODO: Definiere eine Funktion 'schnellangriff'
// - Die Funktion akzeptiert einen Parameter, ob der Wächter wach ist
// - Die Funktion gibt 'true' zurück, wenn ein Schnellangriff gemacht werden 
//	 kann, weil der Ritter erst seine Rüstung anziehen muss wenn er aufsteht

// TODO: Definiere eine Funktion 'spionieren'
// - Die Funktion akzeptiert drei Parameter, ob der Ritter, der Schütze, oder der Prinz wach ist
// - Die Funktion gibt nur 'true' zurück, wenn mindestens einer Wach ist. Weil sonst gibt es nichts zu sehen.

// TODO: Definiere eine Funktion 'brieftaube'
// - Die Funktion akzeptiert zwei Parameter, ob der Schütze, oder der Prinz wach ist
// - Die Funktion gibt nur 'true' zurück, wenn der Schütze schläft und der Prinz wach ist

// TODO: Definiere eine Funktion 'retten'
// - Die Funktion akzeptiert drei Parameter, ob der Ritter, der Schütze, oder der Prinz wach ist
// - Die Funktion gibt nur 'true' zurück, wenn der Prinz wach ist und alle anderen schlafen
