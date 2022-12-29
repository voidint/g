$release="1.4.0"
$os="windows"
$arch="amd64"
$base_dir="$HOME/.g"
$dest_file="${base_dir}/downloads/g${release}.${os}-${arch}.zip"
$url="https://github.com/voidint/g/releases/download/v${release}/g${release}.${os}-${arch}.zip"

function NewDirs () {
    New-Item -Path "$base_dir/downloads","$base_dir/bin" -ItemType "directory"
}

function CleanDirs() {
    Remove-Item -Recurse -Path "$base_dir"
}

function DownloadRelease() {
    Invoke-WebRequest -Uri "$url" -OutFile "$dest_file"
}

function InstallG () {
    Expand-Archive "$dest_file" "$base_dir/bin/"
}

function SetEnv () {
    $e1='$env:GOROOT="$HOME\.g\go"'
    $e2='$env:Path=-join("$HOME\.g\bin;", "$env:GOROOT\bin;", "$env:Path")'

    Out-File -InputObject "" -Append -NoClobber -FilePath $PROFILE
    Out-File -InputObject $e1 -Append -NoClobber -FilePath $PROFILE
    Out-File -InputObject $e2 -Append -NoClobber -FilePath $PROFILE
}


Write-Output "[1/3] Downloading ${url}"
NewDirs
DownloadRelease

Write-Output "[2/3] Install g to the ${HOME}/.g/bin"
InstallG

Write-Output "[3/3] Set environment variables"
SetEnv

Write-Output "Done!"