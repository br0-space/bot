package interfaces

type FortuneInterface interface {
	File() string
	ToMarkdown() string
}

type FortuneServiceInterface interface {
	GetList() []string
	Exists(fileToSearch string) bool
	GetRandomFortune() (FortuneInterface, error)
	GetFortune(file string) (FortuneInterface, error)
}
