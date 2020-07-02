package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//UserSession содержит все данные для авторизации пользователя
type UserSession struct {
	APIKey      string `json:"API_KEY"`
	Destination string `json:"Destination"`
}

//Init парсит файл config.json
func (p *UserSession) Init() error {
	//В случае отсутствия файла - кидаем ошибку
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		return errors.New("File 'config.json' doesnt exist (UserSession.Init())")
	}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return errors.New("Error reading file (UserSession.Init())")
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return errors.New("Invalid config (UserSession.Init())")
	}
	return nil
}

//TestConnection пытается получить на основе информации UserSession данные.
//Функция пытается подключится через http://api_key:<p.APIKey>@<Destination>/api/dashboards/home/
//Возвращает ошибку, если ответ не совпал или APIKey не задан
func (p *UserSession) TestConnection() error {
	response, err := http.Get("http://api_key:" + p.APIKey + "@" + p.Destination + "/api/dashboards/home/")
	if err != nil {
		return errors.New("Response error (UserSession.TestConnection())")
	}
	output, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("Reading responce error (UserSession.TestConnection())")
	}
	var answer map[string]json.RawMessage
	err = json.Unmarshal(output, &answer)
	if err != nil {
		return err
	}
	if answer["message"] != nil {
		return errors.New(string(answer["message"]) + string(" (UserSession.TestConnection())"))
	}
	return nil
}

//GetUIDList получает список UID всех дашбордов, отсекая всё, что не dash-db
func (p *UserSession) GetUIDList() ([]string, error) {
	err := p.TestConnection()
	if err != nil {
		return nil, errors.New("Connection error (UserSession.GetUIDList())")
	}

	//В хорошем случае мы получаем json структуру, которая содержит всю иинформацию о дашбордах и папках
	response, err := http.Get("http://api_key:" + p.APIKey + "@" + p.Destination + "/api/search?folderIds=0&query=&starred=false")
	if err != nil {
		return nil, errors.New("Response error (UserSession.GetIUDList())")
	}
	output, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("Error reading JSON responce (UserSession.GetUIDList())")
	}
	var (
		mappedAnswer []map[string]interface{}
		ret          []string
	)
	err = json.Unmarshal(output, &mappedAnswer)
	if err != nil {
		fmt.Println(string(output))
		return nil, errors.New("Incorrect json answer (UserSession.GetUIDList())")
	}
	for i := range mappedAnswer {
		if mappedAnswer[i]["type"] == "dash-db" {
			ret = append(ret, mappedAnswer[i]["uid"].(string))
		}
	}
	if len(ret) == 0 {
		return nil, errors.New("Dashboards not found (UserSession.GetUIDList())")
	}
	return ret, nil
}
