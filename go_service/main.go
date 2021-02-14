package main

import (
	"os"

	"github.com/andrewawni/chatsystem/app"
)

func main() {
	app := app.App{}
	app.Init(os.Getenv("RABBITMQ_URL"), os.Getenv("REDIS_HOST"))
	app.Run(":8000")
}
