package internal

type Launcher interface {
	Start() error
	Stop() error
}
