package httpserver

import (
	"io"
	"net/http"
)

type Driver struct {
	BaseURL string
	Client  *http.Client
}

func (d Driver) Curse(name string) (string, error) {
	return d.getAndReadFrom(cursePath, name)
}

func (d Driver) Greet(name string) (string, error) {
	return d.getAndReadFrom(greetPath, name)
}

func (d Driver) getAndReadFrom(path string, name string) (string, error) {
	res, err := d.Client.Get(d.BaseURL + path + "?name=" + name)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
