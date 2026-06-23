# Code-signing helper for Albion Market Data Client.
#
# IMPORTANT: A code-signing certificate must be issued to YOUR organization
# ("Albion Market"). You cannot legitimately sign as "Albion Online" /
# "Sandbox Interactive" — that is a different company, and impersonating another
# publisher in a signature is fraud and will be revoked. The name shown in the
# Windows UAC / SmartScreen prompt is the legal subject (CN/O) of the cert.
#
# Options:
#   1. Real publisher trust (no warnings): buy an OV or EV code-signing cert from
#      a CA (DigiCert, Sectigo, SSL.com, ...) issued to "Albion Market". EV gives
#      instant SmartScreen reputation; OV builds reputation over time/downloads.
#   2. Self-signed (testing only): shows "Albion Market" but Windows still warns
#      it's from an unknown publisher. Generate with New-SelfSignedCertificate.
#
# Usage:
#   ./scripts/sign.ps1 -Pfx path\to\albionmarket.pfx -Password 'secret'
#   ./scripts/sign.ps1 -Thumbprint <cert-thumbprint>   # cert already in store

param(
  [string]$Pfx,
  [string]$Password,
  [string]$Thumbprint,
  [string]$Exe = "build/bin/AlbionMarketDataClient.exe",
  [string]$TimestampUrl = "http://timestamp.digicert.com"
)

$ErrorActionPreference = "Stop"
$signtool = (Get-Command signtool.exe -ErrorAction SilentlyContinue)?.Source
if (-not $signtool) {
  # try the Windows SDK default locations
  $signtool = Get-ChildItem "C:\Program Files (x86)\Windows Kits\10\bin" -Recurse -Filter signtool.exe -ErrorAction SilentlyContinue |
    Where-Object { $_.FullName -match "x64" } | Select-Object -First 1 -ExpandProperty FullName
}
if (-not $signtool) { throw "signtool.exe not found. Install the Windows 10/11 SDK." }

if ($Thumbprint) {
  & $signtool sign /sha1 $Thumbprint /fd SHA256 /tr $TimestampUrl /td SHA256 $Exe
} elseif ($Pfx) {
  & $signtool sign /f $Pfx /p $Password /fd SHA256 /tr $TimestampUrl /td SHA256 $Exe
} else {
  throw "Provide either -Thumbprint or -Pfx/-Password."
}

& $signtool verify /pa /v $Exe
Write-Host "Signed and verified: $Exe"
