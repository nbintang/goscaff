param(
  [string]$EnvFile = ".env.local"
)

if (!(Test-Path $EnvFile)) {
  Write-Host "Env file not found: $EnvFile"
  Write-Host "Using Existing Environment Variables..."
  exit 1
}

Get-Content $EnvFile | ForEach-Object {
  $line = $_.Trim()
  if ($line -eq "" -or $line.StartsWith("#")) { return }

  $kv = $line -split "=", 2
  if ($kv.Length -ne 2) { return }

  $key = $kv[0].Trim()
  $val = $kv[1].Trim().Trim('"')

  [Environment]::SetEnvironmentVariable($key, $val, "Process")
}
