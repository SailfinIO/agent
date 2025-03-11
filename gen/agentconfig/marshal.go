package agentconfig

import (
	"bytes"
	"fmt"
)

// Marshal converts the AgentConfig struct into a PKL-formatted []byte.
func Marshal(cfg *AgentConfig) ([]byte, error) {
	var buf bytes.Buffer

	// Optionally include a header if needed:
	buf.WriteString("amends \"./AgentConfig.schema.pkl\"\n\n")

	// Write the serverAddress assignment.
	buf.WriteString(fmt.Sprintf("serverAddress = %q\n", cfg.ServerAddress))

	// Write the apiKey assignment.
	buf.WriteString(fmt.Sprintf("apiKey = %q\n", cfg.ApiKey))

	// Write the remoteHosts list in the expected PKL format.
	buf.WriteString("remoteHosts = List(\n")
	for _, r := range cfg.RemoteHosts {
		buf.WriteString("  (RemoteHost) {\n")
		buf.WriteString(fmt.Sprintf("    host = %q\n", r.Host))
		buf.WriteString(fmt.Sprintf("    user = %q\n", r.User))
		if r.Password != nil {
			buf.WriteString(fmt.Sprintf("    password = %q\n", *r.Password))
		}
		if r.PrivateKey != nil {
			buf.WriteString(fmt.Sprintf("    privateKey = %q\n", *r.PrivateKey))
		}
		buf.WriteString("  },\n")
	}
	buf.WriteString(")\n")

	return buf.Bytes(), nil
}
