package auth

import "net/http"

const (
	APIEndpoint = "/apis/admin.cluster.caicloud.io/v2alpha1/machines"

	inventory = `
{{range $cluster, $value := .outMap}}
[{{$cluster}}]
{{range $index, $node_auth := $value}}
{{$node_auth.host}} ansible_ssh_port={{.sshport}} ansible_ssh_user={{$node_auth.user}} ansible_ssh_pass={{$node_auth.password}}{{end}}
{{end}}
`
)

type config struct {
	masterIP string
	username string
	password string
	sshport  string
}
type client struct {
	httpClient *http.Client
	config config
}