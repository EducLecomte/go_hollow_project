package main

import "github.com/EducLecomte/go_hollow_project/internal/app"

func main() {
	e := app.NewEditorApp()
	if err := e.App.SetRoot(e.Pages, true).Run(); err != nil {
		panic(err)
	}
}
