package main

import (
	"encoding/json"
	"io/ioutil"
)

func main() {
	data, _ := ioutil.ReadFile("./Backups/000000002_backup.json")
	var mappedData = make(map[string]interface{})
	_ = json.Unmarshal(data, &mappedData)
	globalMap := mappedData["dashboard"].(map[string]interface{})
	globalMap = recursiveReplaceTag(globalMap, "editable", false)
	data, _ = json.MarshalIndent(globalMap, "", "\t")
	_ = ioutil.WriteFile("./Changed/000000002.json", data, 0777)
}
