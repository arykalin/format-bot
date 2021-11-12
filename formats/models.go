package formats

type Tag string

type data struct {
	Formats   []Format   `json:"formats"`
	Questions []Question `json:"questions"`
}

type Format struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        []Tag  `json:"tags"`
}

type Question struct {
	Number   int      `json:"number"`
	Question string   `json:"question"`
	Answers  []Answer `json:"answers"`
}

type Answer struct {
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}
