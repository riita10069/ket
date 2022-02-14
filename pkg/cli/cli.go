package cli

type CLI interface {
	Name() string
	Version() string
	Path() string
	Dir() string
	URL() string
	Envs() []string
}
