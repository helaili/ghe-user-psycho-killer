rm -rf release
mkdir release

OS_LIST=(windows darwin linux)

for os in "${OS_LIST[@]}"
do
  echo "Building for $os"
  mkdir release/$os
  env GOOS=$os GOARCH=amd64 go build -ldflags="-s -w" -o release/$os/ghe-user-psycho-killer
  echo "Generating archive release/ghe-user-psycho-killer-$os-amd64.zip"
  if [ $os = "windows" ]
  then
    mv release/$os/ghe-user-psycho-killer release/$os/ghe-user-psycho-killer.exe
    tar -czf "release/ghe-user-psycho-killer-$os-amd64.zip" -C release/$os/ .
  else
    tar -cf "release/ghe-user-psycho-killer-$os-amd64.zip" -C release/$os/ .
  fi
  rm -rf release/$os
done
