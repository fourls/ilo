package provide

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type UnmarshalFunc[T any] func(io.Reader) (*T, error)
type MarshalFunc[T any] func(io.Writer, *T) error

var _ MarshalFunc[int] = YamlMarshal[int]
var _ UnmarshalFunc[int] = YamlUnmarshal[int]

type Loader[T any] interface {
	Load(name string, unmarshalFunc UnmarshalFunc[T]) (*T, error)
}

type Saver[T any] interface {
	Save(name string, value *T, marshalFunc MarshalFunc[T]) error
}

type Provider[T any] interface {
	Loader[T]
	Saver[T]
}

func YamlUnmarshal[T any](reader io.Reader) (*T, error) {
	var data T

	if err := yaml.NewDecoder(reader).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func YamlMarshal[T any](writer io.Writer, value *T) error {
	return yaml.NewEncoder(writer).Encode(value)
}

func NewFileProvider[T any](path string) Provider[T] {
	return fileProvider[T]{basePath: path}
}

func NewConfigProvider[T any]() Provider[T] {
	configPath, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return NewFileProvider[T](filepath.Join(configPath, "ilo"))
}

type fileProvider[T any] struct {
	basePath string
}

func (p fileProvider[T]) ensureBasePath() error {
	return os.MkdirAll(p.basePath, os.ModePerm)
}

func (p fileProvider[T]) makePath(name string) string {
	return filepath.Join(p.basePath, name+".yml")
}

func (p fileProvider[T]) Load(name string, unmarshalFunc UnmarshalFunc[T]) (*T, error) {
	if err := p.ensureBasePath(); err != nil {
		return nil, err
	}

	reader, err := os.Open(p.makePath(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			var def T
			return &def, nil
		} else {
			return nil, err
		}
	}
	defer reader.Close()
	return unmarshalFunc(reader)
}

func (p fileProvider[T]) Save(name string, value *T, marshalFunc MarshalFunc[T]) error {
	if err := p.ensureBasePath(); err != nil {
		return err
	}

	writer, err := os.OpenFile(
		p.makePath(name),
		os.O_CREATE|os.O_WRONLY,
		os.ModePerm)
	if err != nil {
		return err
	}
	defer writer.Close()
	return marshalFunc(writer, value)
}
