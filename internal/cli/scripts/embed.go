package install

import "embed"

//go:embed install-cli.sh
var InstallScript embed.FS
