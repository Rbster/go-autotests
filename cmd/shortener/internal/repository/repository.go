package repository

type UrlInfo struct {
	id       string // ?
	FullUrl  string
	ShortURL string
}

type UrlRepository interface {
	SaveUrlInfo(urlInfo UrlInfo) error
	GetUrlInfoByShortUrl(shortUrl string) (UrlInfo, error)
}

type InMemUrlRepo struct {
	storage map[string]UrlInfo
}

func NewInMemUrlRepo() *InMemUrlRepo {
	storage := make(map[string]UrlInfo)
	return &InMemUrlRepo{
		storage: storage,
	}
}

func (r *InMemUrlRepo) SaveUrlInfo(urlInfo UrlInfo) bool {
	r.storage[urlInfo.ShortURL] = urlInfo
	return true
}

func (r *InMemUrlRepo) GetUrlInfoByShortUrl(shortUrl string) (UrlInfo, bool) {
	url, ok := r.storage[shortUrl]
	return url, ok
}
