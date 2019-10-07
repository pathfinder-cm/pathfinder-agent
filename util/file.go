package util

import (
	"bytes"
	"text/template"

	"github.com/pathfinder-cm/pathfinder-agent/config"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
)

// TODO: to be abstracted
func GenerateBootstrapScriptContent(bs pfmodel.Bootstrapper) (string, int, error) {
	var tpl bytes.Buffer
	var mode int
	if bs.Type == "chef-solo" {
		const content = `
cd /tmp && curl -LO {{.ChefInstaller}} && sudo bash ./install.sh -v {{.ChefVersion}} && rm install.sh
cat > solo.rb << EOF
root = File.absolute_path(File.dirname(__FILE__))
cookbook_path root + "/cookbooks"
EOF
chef-solo -c ~/tmp/solo.rb -j {{.BootstrapAttributes}} {{.CookbooksURL}}
`
		tmpl := template.Must(template.New("content").Parse(content))
		err := tmpl.Execute(&tpl, struct {
			ChefInstaller       string
			ChefVersion         string
			BootstrapAttributes string
			CookbooksURL        string
		}{
			ChefInstaller:       config.ChefInstaller,
			ChefVersion:         config.ChefVersion,
			BootstrapAttributes: bs.Attributes,
			CookbooksURL:        bs.CookbooksUrl,
		})

		if err != nil {
			return "", 0, err
		}

		mode = 600
	}

	return tpl.String(), mode, nil
}
