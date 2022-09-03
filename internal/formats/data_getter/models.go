package data_getter

type Tag string

type data struct {
	Formats   []Format   `json:"formats"`
	Questions []Question `json:"questions"`
}

type Format struct {
	CreatorEmail string `json:"creatorEmail"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Tags         []Tag  `json:"tags"`
}

type QuestionType uint8

const (
	QuestionTypeTag QuestionType = iota
	QuestionTypeLocation
	QuestionTypeAge
)

type Question struct {
	Number   int    `json:"number"`
	Question string `json:"question"`
	Type     QuestionType
	Answers  []Answer `json:"answers"`
}

type Answer struct {
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}
