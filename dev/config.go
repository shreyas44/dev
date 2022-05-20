package dev

type Service struct {
	Cmd string            `json:"cmd"`
	Env map[string]string `json:"env"`
}

type Child struct {
	Path string            `json:"path"`
	Env  map[string]string `json:"env"`
}

type Config struct {
	Init     string             `json:"init"`
	Env      map[string]string  `json:"env"`
	Services map[string]Service `json:"services"`
	Children []Child            `json:"children"`
}
