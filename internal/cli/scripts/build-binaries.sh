# https://akrabat.com/building-go-binaries-for-different-platforms/
# https://akrabat.com/setting-the-version-of-a-go-application-when-building/

version=`git describe --tags HEAD`

platforms=(
"darwin/amd64"
"darwin/arm64"
"linux/amd64"
"linux/arm"
"linux/arm64"
"windows/amd64"
)

for platform in "${platforms[@]}"
do
    # 1
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    os=$GOOS
    if [ $os = "darwin" ]; then
        os="macOS"
    fi

    output_name="deeploy-${version}-${os}-${GOARCH}"
    if [ $os = "windows" ]; then
        output_name+='.exe'
    fi
    
    # 2
    echo "Building release/$output_name..."
    env GOOS=$GOOS GOARCH=$GOARCH go build \
      -C cmd/cli \
      -ldflags "-X github.com/akrabat/deeploy/commands.Version=$version" \
      -o release/$output_name
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting.'
        exit 1
    fi

    # 3
    zip_name="deeploy-${version}-${os}-${GOARCH}"
    pushd cmd/cli/release > /dev/null
    if [ $os = "windows" ]; then
        zip $zip_name.zip $output_name
        rm $output_name
    else
        chmod a+x $output_name
        gzip $output_name
    fi
    popd > /dev/null
done
