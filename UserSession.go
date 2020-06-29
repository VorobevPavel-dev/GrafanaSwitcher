package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"unicode/utf8"
)

//UserSession contains all data from ./JSON/config.json
type UserSession struct {
	APIKey      string `json:"API_KEY"`
	Destination string `json:"Destination"`
	OutputFile  string `json:"Output"`
}

//Init - config parsing
func (p *UserSession) Init() error {
	//Case config.json does not exist
	if _, err := os.Stat("./JSON/config.json"); os.IsNotExist(err) {
		return errors.New("File 'config.json' doesnt exist (UserSession.Init())")
	}
	//Another case
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

//GetDashboardModel try to get JSONModel for dashboard with specific UID
func (p *UserSession) GetDashboardModel(uid string) ([]byte, error) {
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
	jsonString = repairJSON(jsonString)
	err = ioutil.WriteFile(uid+".json", jsonString, 0777)
	err = ioutil.WriteFile("./Backups/"+uid+"_backup.json", jsonString, 0777)
	return output, nil
}

//PostDashboardModel sends all data in JSONFile.json to update JSONModel
//Dashboard's UID already must be in json file
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
	_, err = client.Do(req)
	if err != nil {
		return errors.New("Responce error (UserSession.PostDashboardModel())")
	}
	return nil
}
