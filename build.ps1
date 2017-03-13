###Variables
$env:GOOS = "linux"
$env:GOARCH = "arm"
$MyPath = Split-Path $MyInvocation.MyCommand.Definition

Push-Location $MyPath

& go build .

Pop-Location