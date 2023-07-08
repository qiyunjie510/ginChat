package main

import (
	"ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	err := r.Run(":8081")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}
