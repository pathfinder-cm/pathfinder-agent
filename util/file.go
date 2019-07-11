package util

import (
	"fmt"
	"io"
	"os"

	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
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

func SetupBootstrapFileContent(bs pfmodel.Bootstrapper, content string, mode int) (string, int) {
	switch bs.Type {
	case "chef-solo":
		content = `
cd /tmp && curl -LO https://www.chef.io/chef/install.sh && sudo bash ./install.sh -v 14.12.3 && rm install.sh
cat > solo.rb << EOF
root = File.absolute_path(File.dirname(__FILE__))
cookbook_path root + "/cookbooks"
EOF
`
		execChefSoloCmd := fmt.Sprintf("chef-solo -c ~/tmp/solo.rb -j %s %s", bs.Attributes, bs.CookbooksUrl)
		content = content + "\n" + execChefSoloCmd
		mode = 600
	default:
		content = content
		mode = mode
	}

	return content, mode
}
