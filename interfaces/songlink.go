package interfaces

type SonglinkEntryInterface interface {
	ToMarkdown() string
}

type SonglinkServiceInterface interface {
	GetEntryForUrl(url string) (SonglinkEntryInterface, error)
}
