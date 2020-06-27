#!/usr/bin/env pwsh

$BinDir = "${HOME}/.g/go/bin"

if (Test-Path $BinDir -PathType Container) {
    $User = [EnvironmentVariableTarget]::User
    $Path = [Environment]::GetEnvironmentVariable('Path', $User)

    $BinDir = (Resolve-Path -Path  $BinDir).ToString()
    if (!(";$Path;".ToLower() -like "*;$BinDir;*".ToLower())) {
        [Environment]::SetEnvironmentVariable('Path', "$BinDir;$Path", $User)
        $Env:Path += ";$BinDir"
    }
    else {
        "${BinDir} is already in the `$Path` envirment variable."
    }
}
else {
    "You should run `g.exe` and install any version of golang package first."
}
