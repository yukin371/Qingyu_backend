param()

Write-Error @"
This script is deprecated.
Use the Python entrypoint instead:

python .\scripts\e2e_publication_flow.py --base-url http://localhost:9090
"@

exit 1
