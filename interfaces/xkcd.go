package interfaces

type XkcdComicInterface interface {
	Number() int
	ImageURL() string
	ToMarkdown() string
}

type XkcdServiceInterface interface {
	Random() (XkcdComicInterface, error)
	Latest() (XkcdComicInterface, error)
	Comic(number int) (XkcdComicInterface, error)
}
