# Der Compiler

Die Aufgabe ist entworfen, um dir Übung beim Kompilieren und Ausführen von Go Programmen zu geben.

## Ziel

1. Der Umgang mit dem **Judge** ist bekannt. Du kannst deine Aufgabe herunterladen und wenn du fertig bist mit der Bearbeitung wieder hochladen.
2. Wenn dir jemand seinen **Go Code** schickt weißt du wie es auf deinem Computer **ausgeführt werden** kann


## Code

1. Der Ablauf ist immer gleich: 
    - Downloade den Code
    - Mache Anpassungen am Code
    - Lade den Code wieder hoch und der Judge wird es benoten

2. Für das Herunterladen und Hochladen deines Codes gibt es ein Programm zur Kommunikation mit dem Judge
    - Zum installieren von diesem Programm gib folgenden Befehl in die Eingabeaufforderung:  
    ```
    go install github.com@todo: cli
    ```
    - Wenn du jetzt den Befehl `judge` eingibst, solltest du eine Nachricht bekommen
    todo: bild

3. Im Folgenden ist eine Vorschau vom Code für diese Aufgabe

{{Code}}

4. Downloade den Code mit:   
`judge download 01`


## Schritte

1. Öffne die Eingabeaufforderung

![img](/tasks/01/image.png) TODO: Besseres Bild

2. Navigiere zu deinem `vhs` Ordner. Dazu hier ein Beispiel, um zwischen Ordnern zu navigieren:

Schauen welche Inhalte der Ordner hat
```powershell
C:\Users\Julien> ls
Desktop Bilder Videos
```

Den nächsten Ordner auswählen
```powershell
C:\Users\Julien> cd Desktop
```

Nochmal schauen
```powershell
C:\Users\Julien\Desktop> ls
vhs 
```

Wiederholen bis man am richtigen Ort ist
```powershell
C:\Users\Julien\Desktop> cd vhs
C:\Users\Julien\Desktop\vhs> cd 01
C:\Users\Julien\Desktop\vhs\01> go run main.go
```

