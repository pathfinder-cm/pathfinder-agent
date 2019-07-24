package util

import (
	"io"
	"os"
	"text/template"

	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
	"github.com/pathfinder-cm/pathfinder-agent/config"
)

func WriteStringToFile(filename string, data string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return file, err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return file, err
	}
	return file, nil
}

func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func GenerateBootstrapFileContent(bs pfmodel.Bootstrapper) (string, int, error) {
	var tmpl string
	var mode int
	if bs.Type == "chef-solo" {
		const content = `
cd /tmp && curl -LO {{.ChefInstaller}} && sudo bash ./install.sh -v {{.ChefVersion}} && rm install.sh
cat > solo.rb << EOF
root = File.absolute_path(File.dirname(__FILE__))
cookbook_path root + "/cookbooks"
EOF
chef-solo -c ~/tmp/solo.rb -j {{.BootstrapAttributes}} {{.CookbooksUrl}}
`
		tmpl := template.Must(template.New("content").Parse(content))
		err := tmpl.Execute(os.Stdout, struct {
			ChefInstaller       string
			ChefVersion         string
			BootstrapAttributes string
			CookbooksUrl        string
		}{
			ChefInstaller:       config.ChefInstaller,
			ChefVersion:         config.ChefVersion,
			BootstrapAttributes: bs.Attributes,
			CookbooksUrl:        bs.CookbooksUrl,
		})

		if err != nil {
			return "", 0, err
		}

		mode = 600
	}

	return tmpl, mode, nil
}
