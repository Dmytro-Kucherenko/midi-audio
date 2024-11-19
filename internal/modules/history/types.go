package history

type Schema struct {
	Status   Status   `json:"status"`
	Messages []string `json:"messages"`
	Ports    []string `json:"ports"`
	Apps     []string `json:"apps"`
}
