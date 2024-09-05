package config

type User struct {
	AccessPaths map[string]bool `json:"access_paths"`
	Admin       bool            `json:"admin"`
	UserName    string          `json:"username"`
	Password    string          `json:"password"`
}

type Config struct {
	Dir                string          `json:"dir" default:"uploads"`
	Port               string          `json:"port" default:"8080"`
	Depth              string          `json:"depth" default:"3"`
	ConfFile           string          `json:"confFile" default:"./config.json"`
	AllowedHeaderTypes map[string]bool `json:"allowedHeaderTypes"`
	Users              []User          `json:"users"`
}
