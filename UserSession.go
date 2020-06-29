package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"unicode/utf8"
)

//UserSession позволяет
type UserSession struct {
	APIKey      string `json:"API_KEY"`
	Destination string `json:"Destination"`
	OutputFile  string `json:"Output"`
}

//Init - Парсинг конфига
func (p *UserSession) Init() error {
	//Если файл config.json  не существует
	if _, err := os.Stat("./JSON/config.json"); os.IsNotExist(err) {
		return errors.New("File 'config.json' doesnt exist (UserSession.Init())")
	}
	//Если конфиг существует
	data, err := ioutil.ReadFile("./JSON/config.json")
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
//Ответ должен быть отличным от {"message":"Unauthorized"}
//Возвращает ошибку, если ответ не совпал или APIKey не задан
func (p *UserSession) TestConnection() error {
	response, err := http.Get("http://api_key:" + p.APIKey + "@" + p.Destination + "/api/dashboards/home/")
	if err != nil {
		return errors.New("Response error (UserSession.TestConnection())")
	}
	output, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(output))
	if err != nil {
		return errors.New("Reading responce error (UserSession.TestConnection())")
	}
	var answer map[string]json.RawMessage
	err = json.Unmarshal(output, &answer)
	// if err != nil || answer["message"] != nil {
	// 	return errors.New("Cannot recognize")
	// }
	if err != nil {
		return err
	}
	if answer["message"] != nil {
		return errors.New(string(answer["message"]) + string(" (UserSession.TestConnection())"))
	}
	return nil
}

func (p *UserSession) GetDahsboardModel(uid string) ([]byte, error) {
	err := p.TestConnection()
	if err != nil {
		return nil, errors.New("Connection error (UserSession.GetDashboardModel.TestConnection())")
	}
	response, err := http.Get("http://api_key:" + p.APIKey + "@" + p.Destination + "/api/dashboards/uid/" + uid)
	// command := exec.Command("curl", "http://api_key:"+p.APIKey+"@"+p.Destination+"/api/dashboards/uid/"+uid)
	if err != nil {
		return nil, errors.New("Response error (UserSession.GetDashboardModel())")
	}
	output, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("Reading responce error(UserSession.GetDashboardModel())")
	}
	var answer map[string]interface{}
	err = json.Unmarshal(output, &answer)
	// if err != nil || answer["message"] != nil {
	// 	return errors.New("Cannot recognize")
	// }
	if err != nil {
		return nil, err
	}
	if answer["message"] != nil {
		return nil, errors.New(string(answer["message"].(string)) + string(" (UserSession.TestConnection())"))
	}
	parsed := removeMetaTag(answer)
	jsonString, err := json.MarshalIndent(parsed, "", "\t")
	err = ioutil.WriteFile(p.OutputFile, jsonString, 0777)
	return output, nil
}

func (p *UserSession) PostDashboardModel(JSONFile string) error {
	//Для поста нужно удалить meta тэг
	err := p.TestConnection()
	if err != nil {
		return errors.New("Connection error (UserSession.PostDashboardModel())")
	}
	url := "http://" + string(p.Destination) + "/api/dashboards/db"
	newValues, err := ioutil.ReadFile(JSONFile)
	if err != nil {
		return errors.New("Error reading file(UserSession.PostDashboardModel())")
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(newValues))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", string(utf8.RuneCountInString(string(newValues))))
	req.Header.Set("Authorization", "Bearer "+string(p.APIKey))

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return errors.New("Responce error (UserSession.PostDashboardModel())")
	}
	output, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("Responce reading error (UserSession.PostDashboardModel())")
	}
	fmt.Println(string(output))
	// command := exec.Command("curl",
	// 	"-X", "POST",
	// 	"--insecure",
	// 	"-d", "@"+string(p.OutputFile),
	// 	"-H", "\"Content-type: application/json\"",
	// 	"-H", "\"Authorization: Bearer "+string(p.APIKey)+"\"",
	// 	"http://"+string(p.Destination)+"/api/dashboards/db")
	// output, err := command.Output()
	// if err != nil {
	// 	return err
	// }
	// if output != nil {
	// 	fmt.Println(string(output))
	// }
	// return nil
	// data, _ := ioutil.ReadFile("output.json")
	// temp := string(data)
	// fmt.Println(temp)
	// command := exec.Command(
	// 	"curl",
	// 	"-X", "POST",
	// 	"-d", temp,
	// 	"-H", "\"Content-type: application/json\"",
	// 	"http://api_key:eyJrIjoiYWJnRDdJMXNUTGNlbG5rNjRkVkp1ZVUwaUd2QWdlMWkiLCJuIjoic3R1ZGVudCIsImlkIjoxfQ==@carbon-view.unix.tensor.ru/api/dashboards/db")
	// fmt.Println(command.)
	// 	output, err := command.Output()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(output))
	return nil
}
