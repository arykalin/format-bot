package data_getter

type getter struct {
}

type Getter interface {
	GetData() (formats []Format, questions []Question, err error)
}

func (g *getter) GetData() (formats []Format, questions []Question, err error) {

}

func NewGetter() *getter {
	return &getter{}
}
