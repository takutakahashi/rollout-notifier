module github.com/takutakahashi/rollout-notifier

go 1.15

require (
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/labstack/gommon v0.3.0
	github.com/slack-go/slack v0.7.2
	github.com/urfave/cli v1.22.5
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29 // indirect
	sigs.k8s.io/controller-runtime v0.6.4
)
