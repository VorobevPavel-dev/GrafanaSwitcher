package main

func removeMetaTag(temp map[string]interface{}) map[string]interface{}{
	ret := make(map[string]interface{})
	for key,value := range(temp){
		if (key != "meta"){
			ret[key]=value
		}
	}
	return ret
}