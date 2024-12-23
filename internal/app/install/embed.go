package install

import "embed"

//go:embed install.sh
var InstallScript embed.FS
