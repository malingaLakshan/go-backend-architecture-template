Get-ChildItem -Path . -Recurse -File |
    Select-String -Pattern 'Configuration:'