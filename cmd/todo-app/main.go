package main

import "todo-app/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir, "main")
}
