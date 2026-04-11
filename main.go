package main

func main() {
	e := NewEditorApp()
	if err := e.App.SetRoot(e.Pages, true).Run(); err != nil {
		panic(err)
	}
}
