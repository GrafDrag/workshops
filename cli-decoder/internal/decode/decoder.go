package decode

type Decoder interface {
	Save(str string) error
	validator() bool
	getMD5() string
}
