// Code generated from Pkl module `SailfinIO.agent.AgentConfig`. DO NOT EDIT.
package agentconfig

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type AgentConfig struct {
	// The address the agent listens on.
	ServerAddress string `pkl:"serverAddress"`

	// API key for authenticating agent operations.
	ApiKey string `pkl:"apiKey"`

	// Configuration for remote hosts.
	RemoteHosts []*RemoteHost `pkl:"remoteHosts"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a AgentConfig
func LoadFromPath(ctx context.Context, path string) (ret *AgentConfig, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := evaluator.Close()
		if err == nil {
			err = cerr
		}
	}()
	ret, err = Load(ctx, evaluator, pkl.FileSource(path))
	return ret, err
}

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a AgentConfig
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*AgentConfig, error) {
	var ret AgentConfig
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
