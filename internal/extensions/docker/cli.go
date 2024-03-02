package docker

import (
	"errors"

	"github.com/docker/cli/cli/command"
)

var (
	cli *command.DockerCli
)

func Init() error {
	if cli != nil {
		return nil
	}
	var err error
	if cli, err = command.NewDockerCli(); err != nil {
		return err
	}
	if cli == nil {
		return errors.New("created invalid docker client")
	}
	return nil

}
