$userPath = [Environment]::GetEnvironmentVariable("Path", "User")

if (($userPath -split ";") -notcontains $exeDir) {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$exeDir", "User")
}

$env:Path = "$env:Path;$exeDir"