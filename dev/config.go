package dev

type Service struct {
	Cmd string            `json:"cmd"`
	Env map[string]string `json:"env"`
}

type Config struct {
	Init     string             `json:"init"`
	Services map[string]Service `json:"services"`
	Children []string           `json:"children"`
}
