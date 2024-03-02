package docker

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/Frank-Mayer/yab/internal/lua"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	yaml "gopkg.in/yaml.v3"
)

func Compose(l *lua.LState) int {
	if err := Init(); err != nil {
		l.RaiseError("error initializing docker client: %s", err)
		return 0
	}

	// get the compose file path from the first argument (default to "docker-compose.yml")
	composeFilePath := l.OptString(1, "docker-compose.yml")

	composeFile, err := os.Open(composeFilePath)
	if err != nil {
		l.RaiseError("error opening compose file: %s", err)
		return 0
	}
	defer composeFile.Close()

	service := compose.NewComposeService(cli)
	ctx := context.Background()
	var project *types.Project
	if strings.HasSuffix(composeFilePath, ".json") {
		decoder := json.NewDecoder(composeFile)
		if err := decoder.Decode(&project); err != nil {
			l.RaiseError("error decoding compose file: %s", err)
			return 0
		}
	} else {
		decoder := yaml.NewDecoder(composeFile)
		if err := decoder.Decode(&project); err != nil {
			l.RaiseError("error decoding compose file: %s", err)
			return 0
		}
	}
	if err := service.Up(ctx, project, api.UpOptions{}); err != nil {
		l.RaiseError("error starting compose project: %s", err)
		return 0
	}

	if _, err := service.Wait(ctx, project.Name, api.WaitOptions{}); err != nil {
		l.RaiseError("error waiting for compose project: %s", err)
		return 0
	}

	return 0
}
