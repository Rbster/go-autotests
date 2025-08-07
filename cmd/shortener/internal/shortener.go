package internal

type ShortenerService interface {
	ToShort(full string, host string) (string, error)
	ToFull(short string) (string, error)
}
