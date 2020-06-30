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

//GetDashboardModel пытается получить JSONModel по UID дашборда
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
	if err != nil {
		return nil, err
	}
	if answer["message"] != nil {
		return nil, errors.New(string(answer["message"].(string)) + string(" (UserSession.TestConnection())"))
	}
	parsed := removeMetaTag(answer)
	jsonString, err := json.MarshalIndent(parsed, "", "\t")
	jsonString = repairJSON(jsonString)

	//Создаются две папки - Backups хранит выгруженные копии JSONModel, они не должны использоваться для изменения полей
	//						Changed хранит уже изменённые копии JSONModel, они могут быть выгружены на сервер
	err = ioutil.WriteFile("./Backups/"+uid+"_backup.json", jsonString, 0777)
	// err = ioutil.WriteFile("./Backups/"+uid+"_backup.json", jsonString, 0777)
	return output, nil
}

//PostDashboardModel отправляет JSONModel на сервер
//UID обновляемого дашборда уже записано в соответствующем файле
// func (p *UserSession) PostDashboardModel(uid string) error {
// 	//Для поста нужно удалить meta тэг
// 	err := p.TestConnection()
// 	if err != nil {
// 		return errors.New("Connection error (UserSession.PostDashboardModel())")
// 	}
// 	url := "http://" + string(p.Destination) + "/api/dashboards/db"
// 	newValues, err := ioutil.ReadFile("./Changed/" + string(uid) + ".json")
// 	if err != nil {
// 		return errors.New("File not found. Probably it hasnt been changed (UserSession.PostDashboardModel())")
// 	}
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(newValues))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Content-Length", string(utf8.RuneCountInString(string(newValues))))
// 	req.Header.Set("Authorization", "Bearer "+string(p.APIKey))

// 	client := &http.Client{}
// 	_, err = client.Do(req)
// 	if err != nil {
// 		return errors.New("Responce error (UserSession.PostDashboardModel())")
// 	}
// 	return nil
// }

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

//GetMap получает map[string]interface{} для конкретной JSONModel по UID.
//Также делает бэкап файла
func (p *UserSession) GetMap(uid string) (map[string]interface{}, error) {
	_, err := p.GetDashboardModel(uid)
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadFile("./Changed/" + uid + ".json")
	var mappedData = make(map[string]interface{})
	err = json.Unmarshal(data, &mappedData)
	if err != nil {
		return nil, errors.New("Cannot convert to json (UserSession.GetMap())")
	}
	return mappedData, nil
}

//ChangeTag берёт соответствющий дашборд, делает бэкап в ./Backups и меняет файл.
//Изменённая версия попадает в папку Change с именем <uid>.json.
//Если дашборд уже получен в ./Changed, то берётся локальный файл
//Возвращает количество изменений при удачной попытке
//TODO: нормально описать ошибки
//TODO: удаление файла ./Changed при ошибке
func (p *UserSession) ChangeTag(uid, tagName string, newValue interface{}) (int, error) {
	//Сначала проверяем подключение
	err := p.TestConnection()
	if err != nil {
		return 0, errors.New("Connection error (UserSession.ChangeTag())")
	}
	//Получаем JSONModel в Backup
	_, err = p.GetDashboardModel(uid)
	if err != nil {
		return 0, errors.New("Cannot find dashboard (UserSession.ChangeTag())")
	}
	//Получаем информацию от файла и помещаем в карту
	var mainMap = make(map[string]interface{})
	data, _ := ioutil.ReadFile("./Backups/" + uid + "_backup.json")
	json.Unmarshal(data, &mainMap)

	//Меняем тэг через вспомогательную функцию
	mainMap, err = changeTag(mainMap, tagName, newValue)
	if err != nil {
		return 0, err
	}

	//Записываем изменения в ./Changed
	err = mapToFile("./Changed/"+uid+".json", mainMap)
	if err != nil {
		return 0, err
	}

	//Отправляем изменения на сервер
	err = p.PostDashboardModel("./Changed/" + uid + ".json")
	if err != nil {
		return 0, err
	}
	return len(routes), nil
}

func (p *UserSession) Recover(uid string) error {
	err := p.TestConnection()
	if err != nil {
		return errors.New("Connection error (UserSession.PostDashboardModel())")
	}
	url := "http://" + string(p.Destination) + "/api/dashboards/db"
	// fmt.Println("./Backups/" + string(uid) + "_backup.json")
	newValues, err := ioutil.ReadFile("./Backups/" + string(uid) + "_backup.json")
	if err != nil {
		return errors.New("File not found. Probably it hasnt been changed (UserSession.PostDashboardModel())")
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
