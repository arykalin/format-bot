package formats

type Formats interface {
	GetTags(data_getter.Question) ([]data_getter.Tag, error)
	GetFormats([]data_getter.Tag) ([]data_getter.Format, error)
	GetQuestion(num int) (question *data_getter.Question)
}
