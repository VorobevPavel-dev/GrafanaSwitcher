# GrafanaSwitcher
Небольшая утилита, позволяющая менять тэги внутри JSON файла. Рекомендуется использовать для работы с Grafana Dashboard JSON Model
# Подготовка
В корневом каталоге проекта следует создать config.json

```json
{
    "Destination": "your-grafana.server",
    "API_KEY": "your-personal-api-key"
}
```
Этот файл необходим для работы с сервером (получение, отправка, восстановление дашбордов)

# Установка
```bash
git clone https://github.com/VorobevPavel-dev/GrafanaSwitcher.git
cd GrafanaSwitcher
go build
```

# Флаги командной строки
  - #### -uid
  Уникальный идентификатор дашборда. Необходим для большинства операций
  - #### -tag
  Имя тэга, которое нужно изменить
  - #### -newValue
  Значение, которое будет применено ко всем полям с именем -tag
  - #### -post
  Используется для отправки изменений на сервер
  - #### -uids-only
  Если этот флаг указан, то программа поместит список всех UID в файл с именем UIDList.txt
  - #### -get-only
  Программа только получит копию JSONModel и поместит её в папку ./Backups с именем UID_backups.json
  - #### -restore  
  Используется для отката JSON Model на предыдущюю версию. Для использования необходим заполненный тэг -uid

# Основные структуры
  ## UserSession
  Структура, которая инициализирует пользователя по его API ключу из config.json.
  ### Поля: 
  - APIKey
  - Destination (адрес сервера Grafana)  
  ### Методы:
  - #### UserSession.Init() error
   Позволяет прочитать config.json и поместить все данные в структуру. Необходимо использовать каждый раз при запуске программы
   Возвращает ошибку, если не существует config.json, его невозможно прочитать или он имеет неверный формат
  - #### UserSession.TestConnection() error
   Отправляет GET запрос. Получает ответ в виде JSON. Если тэг "message" не пуст, то невозможно авторизоваться.
   Возвращает ошибку, если невозможно получить ответ, невозможно его прочитать, или "message" != nil
  -  #### UserSession.GetUIDList() ([]string, error) 
   Отправляет GET запрос. Получает ответ в виде JSON со всеми известными структурами на сервере. Из этого выбираются только те, что имеют тип "dash-db"
   Возвращает массив строк - список всех UID "dash-db"
   Возвращает ошибку, если невозможно подключиться (UserSessoin.TestConnection)
   
 ## Dashboard
  Структура, которая инициализирует дашборд по файлу или его UID.
  ### Поля:
  - UID - уникальный ID дашборда
  - ID - простой ID дашборда
  - Version - версия данного JSONModel
  - Backup - адрес локального неизменённого файла
  - Changed - адрес изменённого файла
  - MainMap - представление данного JSONModel файла в виде карты тэгов
  ### Методы:
  - #### Init(uid string)
  Инициализирует дашборд. Если он был уже получен (существует файл ./Backups/uid_backup.json), то заполняет  абсолютно все поля исходя из этих данных.
  Если файл не существует, то прописывает автоматически поля Backup и Changed, даже если файл не будет получен.
  - #### Get(p* UserSession) error
  Все функции, которые работают с API должны получать UserSession как параметр для подключения к серверу.  
  Фунция получает JSONModel с сервера по Dashboard.UID и сохраняет файл в ./Backups/Dashboard.UID_backup.json. Форматирует полученные данные путём удаления тэга "meta". Также заполняет все поля экземпляра Dashboard необходимыми данными.  
  Возвращает ошибку, если невозможно авторизоваться, ответ содержит тэг "message", пришёл некорректный ответ
  - #### Post(p* UserSession) error
  Функция отправляет файл ./Changed/uid.json на сервер.  
  Возвращает ошибку, если файл не существует, невозмонжо подключиться, пришёл некорректный ответ
  - #### ChangeTag(tagName string, newValue interface{})
  Функция берёт данные из файла ./Backups/uid_backup.json, меняет все значения тэга tagName на newValue  
  Пример:  
  ```json
  {
    "name": "test123",
    "panel":
    {
        "name": "qwerty"
    }
  }
  ```
  ```go
  dashboard.ChangeTag("name", "hello")
  ```
  ```json
  {
    "name": "hello",
    "panel":
    {
        "name": "hello"
    }
  }
  ```
  Изменённые данные записывает в файл ./Changed/uid.json
  - #### Restore (p* UserSession) (*Dashboard, error)
  Функция откатывает изменения на один коммит назад.  
  Можно было бы просто отправить файл ./Backups/uid_backup.json, однако в том файле не изменен тэг "version"  
  Restore вытаскивает dashboard.Version, уменьшает его на один и через POST запрос откатывает изменения. Они отображаются на вкладке Versions в Grafana
  
