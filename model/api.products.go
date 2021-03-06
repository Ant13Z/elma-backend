package model

import (
	"../db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Product struct {
	Name string
	Count int
	Description *string
}

func getData(request *http.Request) Product{
	defer request.Body.Close()
	var element Product
	filter := true
	body, _ := ioutil.ReadAll(request.Body)
	err := json.Unmarshal(body, &element)
	if err != nil {
		//если поменяем фатал на паник, то сервер при отправке не числа падать не будет
		log.Fatal(err)
	} else {
		//по идее, не прямая подстановка была создана для того чтобы решать вопросы безопасности
		//но как это реализовано в го, я не особо в курсе, так что фильтруем дополнительно регулярками
		//https://habr.com/en/post/308088/ говорят не полностью безопасно
		regExpString := regexp.MustCompile("[^a-zA-Z0-9а-яА-ЯёЁ]")

		element.Name = regExpString.ReplaceAllString(element.Name, "")
		if len(element.Name) < 0 || len(element.Name) > 255 {
			filter = false
		}

		if element.Description != nil {
			*element.Description = regExpString.ReplaceAllString(*element.Description, "")
		}

		if element.Count < 0 || element.Count > 9999 {
			//проверка на число не нужна, т.к. типизированный язык
			filter = false
		}
		//ситуации когда пользователь не добился этого быть не может, а если кто то шлет запросы напрямую
		//то ему информация и не нужна
		if filter {
			return element
		}
	}
	var elementReturn Product
	return elementReturn
}

func ProductsHandler(response http.ResponseWriter, request *http.Request){
	//нам без разницы откуда к нам пришел запрос. если пользователь знает куда делать запрос, то может это все и
	//руками сделать, т.е. для такого приложения должна быть доп проверка авторизации
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	response.Header().Set("Access-Control-Allow-Credentials", "false")
	response.Header().Set("Access-Control-Max-Age", "86400")
	response.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding")
	response.Header().Set("Accept", "application/json")
	response.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch request.Method {
		case "GET":
			ans, err := db.Connect.Query("select `name`, `count`, `description` from products order by id asc")
			if err != nil {
				panic(err)
			}
			products := make([]*Product, 0)
			for ans.Next() {
				elem := new(Product)
				err := ans.Scan(&elem.Name, &elem.Count, &elem.Description)
				if err != nil {
					log.Fatal(err)
				}
				products = append(products, elem)
			}
			productsJSON, _ := json.Marshal(products)
			fmt.Fprintf(response, string(productsJSON))
		case "POST":
			fmt.Println("POST")
			element := getData(request)
			if len(element.Name) > 0 {
				_, err := db.Connect.Exec(""+
					"insert ignore into products "+
					"(`name`, `count`, `description`) "+
					"values (?, ?, ?)",
					element.Name,
					element.Count,
					element.Description)
				if err != nil {
					panic(err)
				}
				productJSON, _ := json.Marshal(element)
				fmt.Fprintf(response, string(productJSON))
			}
		case "PUT":
			fmt.Println("PUT")
			element := getData(request)
			if len(element.Name) > 0 {
				_, err := db.Connect.Exec("" +
					"update products " +
					"set `count` = ?, `description` = ? " +
					"where `name` = ?",
					element.Count,
					element.Description,
					element.Name)
				if err != nil {
					panic(err)
				}
				productJSON, _ := json.Marshal(element)
				fmt.Fprintf(response, string(productJSON))
			}
		case "DELETE":
			//не требуется реализация, но у нас же REST ful
			fmt.Println("DELETE")
		case "OPTIONS":
			//CORS политика для post\put запросов
	}
}