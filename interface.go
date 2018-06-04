package translate

type ITranslate interface {
	Translate(string, Option) (string, error)
}

type Option struct {
	From Language
	To   Language
}

type Language string
