package cmd

import (
	"fmt"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/templates"
)

// Dependencies holds all dependencies required by the commands.
type Dependencies struct {
	Config          *config.Config
	Logger          logger.Logger
	FS              fs.FileSystem
	TemplateManager templates.TemplateManager
}

// defaultInputReader is a simple implementation of templates.InputReader that uses standard input.
type defaultInputReader struct{}

func (r *defaultInputReader) ReadResponse() (string, error) {
	var response string
	_, err := fmt.Scanln(&response)
	return response, err
}
