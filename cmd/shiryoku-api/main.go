package main

import shiryoku_routers "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers"

func main() {
	router := shiryoku_routers.GetFilledRouter()	

	// FIXME: port from config
	router.Run(":8080")
}
