package internal

type Deployment struct {
	UID   string            `json:"uid"`
	Image string            `json:"image"`
	Ports map[string]string `json:"ports"`
	Env   map[string]string `json:"env"`
}
