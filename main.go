package main

import (
	"github.com/ZeraVision/zn-wallet-manager/api"
	"github.com/joho/godotenv"
)

func main() {

	//* Load your environment variables via whatever method you prefer
	godotenv.Load(".env")

	go api.StartAPI()

	select {}
}
