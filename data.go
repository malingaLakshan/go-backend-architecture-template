$target = (Get-Location).Path
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")

if ($userPath -notlike "*$target*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$target", "User")
}