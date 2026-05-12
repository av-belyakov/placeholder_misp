package dicontainer

import "github.com/av-belyakov/placeholder_misp/commoninterfaces"

// NewDIContainer ленивая инициализация DI контейнера
func NewDIContainer(rootDir string, ch chan commoninterfaces.Messager) *DiContainer {
	return &DiContainer{
		rootDir: rootDir,
		ch:      ch,
	}
}
