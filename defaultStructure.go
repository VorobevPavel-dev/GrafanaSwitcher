package main

type JSONModel struct {
	annotations struct {
		list []struct {
			builtIn    int    `json:"builtIn"`
			datasource string `json:"datasource"`
			enable     bool   `json:"enable"`
			hide       bool   `json:"enable"`
			iconColor  string `json:"iconColor"`
			name       string `json:"name"`
			Type       string `json:"type"`
		} `json:"list"`
	} `json:"annotations"`
	editable     bool `json:"editable"`
	gnetID       int  `json:"gnetId"`
	graphTooltip int  `json:"graphTooltip"`
	hideControls bool `json:"hideControls"`
	id           int  `json:"id"`
	links        []struct {
		asDropdown bool     `json:"asDropdown"`
		icon       string   `json:"icon"`
		tags       []string `json:"tags"`
		title      string   `json:"title"`
		Type       string   `json:"type"`
	} `json:"links"`
	refresh string `json:"refresh"`
	rows    []struct {
		collapse bool   `json:"collapse"`
		height   string `json:"height"`
		panels   []struct {
			cacheTimeout     string   `json:"cacheTimeout"`
			colors           []string `json:"colors"`
			content          string   `json:"content"`
			datasource       int      `json: datasource`
			editable         bool     `json:"editable"`
			Error            bool     `json:"error"`
			format           string   `json:"format"`
			graphID          string   `json:"graphId"`
			hideTimeOverride bool     `json:"hideTimeOverride"`
			id               int      `json:"id"`
			init             struct {
			} `json:""`
		} `json:"panels"`
		repeat          int    `json:"repeat"`
		repeatIteration int    `json:"repeatIteration"`
		repeatRowID     int    `json:"repeatRowId"`
		showTitle       bool   `json:"showTitle"`
		title           string `json:"title"`
		titleSize       string `json:"titleSize"`
	} `json:"rows"`
}
