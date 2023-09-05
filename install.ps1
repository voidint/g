$release = "1.5.2"
$os = "windows"
$arch = "amd64"
$default_base_dir="$HOME\.g"
$dest_file = "${base_dir}\downloads\g${release}.${os}-${arch}.zip"
$url = "https://github.com/voidint/g/releases/download/v${release}/g${release}.${os}-${arch}.zip"

param([string] $base_dir = "$default_base_dir")


function NewDirs () {
    New-Item -Force -Path "$base_dir\downloads", "$base_dir\bin" -ItemType "directory"
}

function CleanDirs() {
    Remove-Item -Recurse -Path "$base_dir"
}

function DownloadRelease() {
    Invoke-WebRequest -Uri "$url" -OutFile "$dest_file"
}

function InstallG () {
    Expand-Archive "$dest_file" "$base_dir\bin\"
}


function setHOME() {
    if ($base_dir -ne $default_base_dir) {
        [System.Environment]::SetEnvironmentVariable("G_EXPERIMENTAL", "true", [System.EnvironmentVariableTarget]::User)
    }
    [System.Environment]::SetEnvironmentVariable("G_HOME", $base_dir, [System.EnvironmentVariableTarget]::User)
    [System.Environment]::SetEnvironmentVariable("GOROOT", "$base_dir\go", [System.EnvironmentVariableTarget]::User)
}


function setPath() {
    $paths = [System.Environment]::GetEnvironmentVariable("PATH", [System.EnvironmentVariableTarget]::User) -split ';'
    $newPaths = @("%G_HOME%\bin", "%GOROOT%\bin", "%GOPATH%\bin")

    foreach ($p in $newPaths) {
        if ($p -in $paths) {
            Write-Output "$p already exists"
            continue
        }

        [System.Environment]::SetEnvironmentVariable(
            "PATH",
            [System.Environment]::GetEnvironmentVariable("PATH", [System.EnvironmentVariableTarget]::User) + "$p;",
            [System.EnvironmentVariableTarget]::User
        )
        Write-Host -ForegroundColor Green "$p appended"
    }
}

function SetEnv () {
    setHOME
    setPath
}

Write-Host -ForegroundColor Blue "[1/3] Downloading ${url}"
NewDirs
DownloadRelease

Write-Host -ForegroundColor Blue "[2/3] Install g to the ${base_dir}\bin"
InstallG

Write-Host -ForegroundColor Blue "[3/3] Set environment variables"
SetEnv

Write-Host -ForegroundColor Green "Done!"
