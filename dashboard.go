package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

//Dashboard содержит все данные про определённый grafana dashboard
//Backup и Changed - поля в которых записан путь до бэкап файла и изменённого файла соответственно
//MainMap - представление json model как карты тэгов
type Dashboard struct {
	UID     string
	ID      int
	Version int
	Backup  string
	Changed string
	MainMap map[string]interface{}
}

//Init производит начальную инициализацию полей Backup и Changed.
//Также при имеющемся файле в папке Backups иинициализирует все поля
func (d *Dashboard) Init(uid string) {
	d.UID = uid
	d.Backup = "./Backups/" + d.UID + "_backup.json"
	d.Changed = "./Changed/" + d.UID + ".json"

	//Если файл уже был получен (есть копия в Backups), то получаем из него карту и ставим как MainMap
	//Также ищем id и version
	if _, err := os.Stat(d.Backup); !os.IsNotExist(err) {
		var data = make(map[string]interface{})
		text, _ := ioutil.ReadFile(d.Backup)
		json.Unmarshal(text, &data)
		d.MainMap = data
		innerMap := d.MainMap["dashboard"].(map[string]interface{})
		tempID := innerMap["id"].(float64)
		d.ID = int(tempID)
		tempVersion := innerMap["version"].(float64)
		d.Version = int(tempVersion)
	}
}

//Get получает dashboard на основе данных из config.json и uid внутри структуры.
//Перезаписывает файл внутри папки Backups
func (d *Dashboard) Get(p *UserSession) error {

	//Проверяем валидность API ключа через тестовое подключение пользователя
	err := p.TestConnection()
	if err != nil {
		return errors.New("Connection error (Dashboard.Get())")
	}

	//Выполняем GET запрос для получения JSON Model
	response, err := http.Get("http://api_key:" + p.APIKey + "@" + p.Destination + "/api/dashboards/uid/" + d.UID)
	if err != nil {
		return errors.New("Response error (Dashboard.Get())")
	}

	//Пытаемся прочитать ответ. В случае несоответствия выкидываем ошибку
	output, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("Reading responce error (Dashboard.Get())")
	}

	//Составляем главную карту для JSON Model
	var answer map[string]interface{}
	err = json.Unmarshal(output, &answer)
	d.MainMap = answer

	//Ищем версию и id для json model
	innerMap := d.MainMap["dashboard"].(map[string]interface{})
	tempID := innerMap["id"].(float64)
	d.ID = int(tempID)
	tempVersion := innerMap["version"].(float64)
	d.Version = int(tempVersion)

	if err != nil {
		return errors.New("Cannot convert responce to json struct (Dashboard.Get())")
	}
	//Если message не пуст, то он несёт в себе ошибку авторизации либо иную ошибку
	if answer["message"] != nil {
		return errors.New(string(answer["message"].(string)) + string(" (Dashboard.Get())"))
	}

	//Удаляем тэг с метаданными
	parsed := removeMetaTag(answer)

	//Приводим карту в читаемый вид
	jsonString, err := json.MarshalIndent(parsed, "", "\t")

	//Заменяем символы &,<,> на нормальные
	jsonString = repairJSON(jsonString)

	//Записываем полученную карту в бэкап файл
	err = ioutil.WriteFile(d.Backup, jsonString, 0777)
	return nil
}

//Post отправляет изменения на сервер.
//UID дашборда, который обновится, уже записан в json файле d.Changed
func (d *Dashboard) Post(p *UserSession) error {
	//Проверяем соединение через UserSession.ApiKey
	err := p.TestConnection()
	if err != nil {
		return errors.New("Connection error (Dashboard.Post())")
	}

	//Формируем POST запрос в соответствии с API
	url := "http://" + string(p.Destination) + "/api/dashboards/db"
	newValues, err := ioutil.ReadFile(d.Changed)
	if err != nil {
		return errors.New("Error reading file(Dashboard.Post())\nProbably there is no Changed copy of this dashboard")
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(newValues))

	//Добавляем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", string(utf8.RuneCountInString(string(newValues))))

	//Необходимый заголовок для авторизации (https://grafana.com/docs/grafana/latest/http_api/dashboard/)
	req.Header.Set("Authorization", "Bearer "+string(p.APIKey))

	//Выполняем POST запрос
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return errors.New("Responce error (Dashboard.Post())")
	}
	d.Version++
	return nil
}

//ChangeTag делает копию файла из бэкапа в папку Changed и меняет в нём тэг.
func (d *Dashboard) ChangeTag(tagName string, newValue interface{}) error {
	//Меняем тэг через вспомогательную функцию (util.go)
	var err error
	d.MainMap, err = changeTag(d.MainMap, tagName, newValue)
	if err != nil {
		return errors.New("Cannot change value in map[string]interface{} (Dashboard.ChangeTag().changeTag()")
	}

	//Записываем изменения в Changed через вспомогательную функцию (util.go)
	err = mapToFile(d.Changed, d.MainMap)
	if err != nil {
		return errors.New("Cannot write file to Change folder (Dashboard.ChangeTag())")
	}
	return nil
}

//Restore позволяет откатить изменения на предыдущую версию.
//Далее выполнить POST /api/dashboards/id/:dashboardId/restore.
//Однако, данные о версии следует передать через JSON файл (https://grafana.com/docs/grafana/latest/http_api/dashboard_versions/)
//Если всё прошло успешно - перезаписываем ./Backups и в любом случае удаляем ./Changed.
//Функция вернёт тот же дашборд, если произошла ошибка. Если всё прошло хорошо - откаченную версию
func (d *Dashboard) Restore(p *UserSession) (*Dashboard, error) {
	//Сначала проверяем подключение
	err := p.TestConnection()
	if err != nil {
		return d, errors.New("Connection error (UserSession.ChangeTag())")
	}

	//Удаляем изменённый файл
	_ = os.Remove(d.Changed)

	//Необходимо составить запрос для восстановления предыдущей версии
	url := "http://" + string(p.Destination) + "/api/dashboards/id/" + strconv.Itoa(d.ID) + "/restore"
	versionJSON := "{\"version\": " + strconv.Itoa(d.Version-1) + "}"
	// fmt.Println(string(versionJSON))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(versionJSON)))

	//Добавляем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	//Необходимый заголовок для авторизации (https://grafana.com/docs/grafana/latest/http_api/dashboard/)
	req.Header.Set("Authorization", "Bearer "+string(p.APIKey))

	//Выполняем запрос
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return d, errors.New("Cannot restore dashboard (Dashboard.Restore())")
	}
	textOutput, _ := ioutil.ReadAll(response.Body)
	//Если в ответе нет success - кидаем ошибку
	if !strings.Contains(string(textOutput), string("success")) {
		return d, errors.New("Critical failure (Dashboard.Restore())\n" + string(textOutput))
	}
	if err != nil {
		return d, errors.New("Responce error (Dashboard.Post())")
	}

	//Так как всё прошло хорошо - удаляем backup и записываем новый
	_ = os.Remove(d.Backup)

	//Обновляем дашборд локально
	newDashboard := new(Dashboard)
	newDashboard.Init(d.UID)
	newDashboard.Get(p)
	return newDashboard, nil
}

//Print - более читаемое представление о дашборде (необходимо для отладки и поиска ошибок)
func (d *Dashboard) Print() {
	fmt.Println("Dashboard", d.UID)
	fmt.Println("\tID:", strconv.Itoa(d.ID))
	fmt.Println("\tVersion:", strconv.Itoa(d.Version))
	fmt.Println("\tBackup file:", d.Backup)
	fmt.Println("\tChanged version:", d.Changed)
}
