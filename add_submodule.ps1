param(
    [string]$submodulePath,
    [string]$submoduleUrl,
    [string]$submoduleName
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

# Set the submodule name if provided, otherwise use the default (path)
if (-not $submoduleName) {
    $submoduleName = $submodulePath
}

# 1. Add the submodule
git submodule add -b main --name $submoduleName $submoduleUrl $submodulePath

# 2. Initialize and update submodules recursively
git submodule update --init --recursive

# 3. Stage the changes
git add .gitmodules $submodulePath

# 4. Commit the changes
git commit -m "Added submodule $submoduleName"

Write-Host "Submodule '$submoduleName' added successfully at '$submodulePath'."