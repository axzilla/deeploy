#!/bin/bash

platforms=(
"darwin/amd64"
"darwin/arm64"
"linux/amd64"
"linux/arm64"
)

filename=""
for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    OS=${platform_split[0]}
    ARCH=${platform_split[1]}
    if [[ $(uname) == "$OS" && $(uname -m) == "$ARCH" ]]; then
        filename="deeploy-${OS}-${ARCH}"
    fi
done

downloadURL="https://github.com/axzillal/deeploy/releases/latest/download/${filename}"
curl -LO "$downloadURL"
chmod +x "$filename"
sudo mv "$filename" /usr/local/bin/deeploy
