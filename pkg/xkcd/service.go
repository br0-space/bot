package xkcd

import (
	"context"
	"github.com/br0-space/bot/interfaces"
	xkcdv2 "github.com/nishanths/go-xkcd/v2"
	"math/rand"
)

type Service struct{}

func MakeService() Service {
	return Service{}
}

func (s Service) Random() (interfaces.XkcdComicInterface, error) {
	latest, err := s.Latest()
	if err != nil {
		return Comic{}, err
	}

	number := rand.Intn(latest.Number())

	return s.Comic(number)
}

func (s Service) Latest() (interfaces.XkcdComicInterface, error) {
	comic, err := s.getClient().Latest(context.Background())
	if err != nil {
		return Comic{}, err
	}

	return FromComic(comic), nil
}

func (s Service) Comic(number int) (interfaces.XkcdComicInterface, error) {
	comic, err := s.getClient().Get(context.Background(), number)
	if err != nil {
		return Comic{}, err
	}

	return FromComic(comic), nil
}

func (s Service) getClient() *xkcdv2.Client {
	return xkcdv2.NewClient()
}
