param(
    [string]$submodulePath,
    [string]$submoduleUrl
)

# Check if the submodule path and URL are provided
if (-not $submodulePath) {
    Write-Error "Please provide the path to the submodule using the -submodulePath parameter."
    exit 1
}
if (-not $submoduleUrl) {
    Write-Error "Please provide the URL of the submodule using the -submoduleUrl parameter."
    exit 1
}

# 1. Add the submodule
git submodule add $submoduleUrl $submodulePath

# 2. Stage the changes
git add .gitmodules $submodulePath

# 3. Commit the changes
git commit -m "Added submodule $submodulePath"

Write-Host "Submodule '$submodulePath' added successfully."