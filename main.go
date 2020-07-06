package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	//Переменная со всеми путями до тэгов. Нужна в методе dashboard.ChangeTag()
	routes []string

	//Restore - флаг для обозначения процедуры восстановления.
	//Если этот флаг указан, то Tag и NewValue просто игнорируются
	Restore bool

	//Post - флаг для обозначения процедуры обновления
	//Если флаг указан изменённая версия дашборда будет отправлена на сервер
	Post bool

	//Get - флаг для обохначения процедуры получения
	//Если флаг поднят, то будет просто получена версия JSONModel в папку ./Backups
	Get bool

	//UIDList - флаг для обозначения процедуры получения списка всех UID на сервере
	//Если флаг поднят, то будет получен список всех UID и добавлен в файл .UIDList.txt
	UIDList bool

	//Tag - название тэга для изменения
	Tag string

	//NewValue - новое значения для тэга
	NewValue string

	//UID - UID дашборда
	UID string
)

func init() {
	//Проверка необходимых директорий и их создание
	if _, err := os.Stat("./Backups"); os.IsNotExist(err) {
		os.Mkdir("./Backups", 0777)
	}
	if _, err := os.Stat("./Changed"); os.IsNotExist(err) {
		os.Mkdir("./Changed", 0777)
	}
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		_, e := os.Create("config.json")
		if e != nil {
			fmt.Println("Cannot create config.json file. Try to create it manually")
			os.Exit(1)
		}
		fmt.Println("File config.json has been created. Check README.md to find out necessary information")
		os.Exit(1)
	}

	//Обработка флагов командной строки
	flag.StringVar(&NewValue, "newValue", "", "New value that replace old values")
	flag.StringVar(&Tag, "tag", "", "Tag name that must be replaced with a -newValue")
	flag.BoolVar(&Restore, "restore", false, "If this flag is raised you will restore your json model from previous version. Flags -newValue and -tag will be ignored")
	flag.StringVar(&UID, "uid", "", "UID of a dashboard that should be changed or restored. Necessary flag")
	flag.BoolVar(&Post, "post", false, "If this flat is raised all changes will be posted on a server")
	flag.BoolVar(&Get, "get-only", false, "If this flag is raised you will get a copy of JSONModel. All other flags and commangs will be ignored")
	flag.BoolVar(&UIDList, "uids-only", false, "If this flag is raised you will get list of all UID on a server. All other flags will be ignored")
	flag.Parse()

	if UID == "" {
		if !Get {
			fmt.Println("-uid must not be empty")
			os.Exit(1)
		}
		if !UIDList {
			fmt.Println("-uid must not be empty")
			fmt.Println("Choose -get or -uids-only flag for this case")
			os.Exit(1)
		}
	}

	if Tag != "" && NewValue == "" {
		fmt.Println("A new value is empty. Are you sure you wand to continue? (y/n)")
		var choise string
		_, _ = fmt.Scan(&choise)
		if choise != "y" {
			fmt.Println("Operation aborted")
			os.Exit(1)
		}
	}

	if Get && UIDList {
		fmt.Println("Incorrect combination of flags. Check -h for help")
		os.Exit(1)
	}
}

func main() {
	user := new(UserSession)
	user.Init()
	dashboard := new(Dashboard)
	if UIDList {
		list, err := user.GetUIDList()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		output := ""
		for i := range list {
			output += list[i] + "\n"
		}
		_ = ioutil.WriteFile("UIDList.txt", []byte(output), 0777)
		os.Exit(0)
	}
	dashboard.Init(UID)
	err := dashboard.Get(user)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//Если Get поднят, то выходим из программы, так как JSONModel мы получили ранее
	if Get {
		os.Exit(0)
	}
	//Если поднят флаг Restore - откатываем изменения
	if Restore {
		_, err = dashboard.Restore(user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	} else {
		err = dashboard.ChangeTag(Tag, NewValue)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if Post {
			err = dashboard.Post(user)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
	os.Exit(0)
}
