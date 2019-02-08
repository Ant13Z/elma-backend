package daemon

import (
	"../db"
	"../model"
	"log"
	"net/http"
)

type Config struct {
	Port string
}

func Run(cfgServer Config, cfgDB db.Config){
	//коннект к бд
	//потенциально можно написать штуку которая будет переподключаться в случае обрыва соединения
	db.Connect = db.ConnectMySQL(cfgDB)
	//роутинг
	http.HandleFunc("/api/products/", model.ProductsHandler)
	//создаем сервер
	err := http.ListenAndServe(cfgServer.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return
}