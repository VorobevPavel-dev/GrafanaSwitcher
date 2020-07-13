package main

import (
	"encoding/json"
	"errors"
)

//Panel is a default panel structure
type Panel struct {
	AliasColors map[string]string `json:"aliasColors"`
	Bars        bool              `json:"bars"`
	DashLength  int               `json:"dashLength"`
	Dashes      bool              `json:"dashes"`
	Datasource  interface{}       `json:"datasource"`
	FieldConfig struct {
		Defaults struct {
			Custom map[string]interface{} `json:"custom"`
		} `json:"defaults"`
		Overrides []interface{} `json:"overrides"`
	} `json:"fieldConfig"`
	Fill         int `json:"fill"`
	FillGradient int `json:"fillGradient"`
	GridPos      struct {
		H int `json:"h"`
		W int `json:"w"`
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"gridPos"`
	HiddenSeries bool `json:"hiddenSeries"`
	ID           int  `json:"id"`
	Legend       struct {
		Avg     bool `json:"avg"`
		Current bool `json:"current"`
		Max     bool `json:"max"`
		Min     bool `json:"min"`
		Show    bool `json:"show"`
		Total   bool `json:"total"`
		Values  bool `json:"values"`
	} `json:"legend"`
	Lines           bool                   `json:"lines"`
	Linewidth       int                    `json:"linewidth"`
	NullPointMode   string                 `json:"nullPointMode"`
	Options         map[string]interface{} `json:"options"`
	Percentage      bool                   `json:"percentage"`
	Pointradius     int                    `json:"pointradius"`
	Points          bool                   `json:"points"`
	Renderer        string                 `json:"renderer"`
	SeriesOverrides []interface{}          `json:"seriesOverrides"`
	SpaceLength     int                    `json:"spaceLength"`
	Stack           bool                   `json:"stack"`
	SteppedLine     bool                   `json:"steppedLine"`
	Targets         []struct {
		RefID  string `json:"refId"`
		Target string `json:"target"`
	} `json:"targets"`
	Thresholds  []interface{} `json:"thresholds"`
	TimeFrom    interface{}   `json:"timeFrom"`
	TimeRegions []interface{} `json:"timeRegions"`
	TimeShift   interface{}   `json:"timeShift"`
	Title       string        `json:"title"`
	Tooltip     struct {
		Shared    bool   `json:"shared"`
		Sort      int    `json:"sort"`
		ValueType string `json:"value_type"`
	} `json:"tooltip"`
	Type  string `json:"type"`
	Xaxis struct {
		Buckets interface{}   `json:"buckets"`
		Mode    string        `json:"mode"`
		Name    interface{}   `json:"name"`
		Show    bool          `json:"show"`
		Values  []interface{} `json:"values"`
	} `json:"xaxis"`
	Yaxes []struct {
		Format  string      `json:"format"`
		Label   interface{} `json:"label"`
		LogBase int         `json:"logBase"`
		Max     interface{} `json:"max"`
		Min     interface{} `json:"min"`
		Show    bool        `json:"show"`
	} `json:"yaxes"`
	Yaxis struct {
		Align      bool        `json:"align"`
		AlignLevel interface{} `json:"alignLevel"`
	} `json:"yaxis"`
}

//ParseJSONString allows to convert JSONString to Panel struct
func (p *Panel) ParseJSONString(JSONData string) error {
	data := []byte(JSONData)
	err := json.Unmarshal(data, &p)
	if err != nil {
		return errors.New("Cannot convert JSONPanel to string (Panel.ParseJSONString())")
	}
	return nil
}

//ConvertToJSON returns a string type of current Panel
func (p *Panel) ConvertToJSON(indent bool) (string, error) {
	var (
		data []byte
		err  error
	)
	if indent {
		data, err = json.MarshalIndent(p, "", "\t")
	} else {
		data, err = json.Marshal(p)
	}
	if err != nil {
		return "", errors.New("Cannot convert to JSON (Panel.ConvertToJSON())")
	}
	return string(data), nil
}

//ChangeTarget returns true if replace
func (p *Panel) ChangeTarget(oldTarget, newTarget string) (*Panel, bool) {
	for i := range p.Targets {
		if p.Targets[i].Target == oldTarget {
			p.Targets[i].Target = newTarget
			return p, true
		}
	}
	return p, false
}
