#!/bin/bash
echo "=== Checking all API files for shared.Success/SuccessData/Error calls ==="
echo ""
echo "P0 Core Modules:"
echo "=== social ==="
find api/v1/social -name "*_api.go" ! -name "*_test.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null
echo "=== reader ==="
find api/v1/reader -name "*_api.go" ! -name "*_test.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null
echo "=== search ==="
find api/v1/search -name "*_api.go" ! -name "*_test.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null
echo ""
echo "P1 Auxiliary Modules:"
echo "=== ai ==="
find api/v1/ai -name "*_api.go" ! -name "*_test.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null
echo "=== stats ==="
find api/v1/stats -name "*_api.go" ! -name "*_test.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null
echo "=== recommendation ==="
find api/v1/recommendation -name "*.go" ! -name "*_test.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null
echo "=== version ==="
find api/v1 -name "version_api.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null
echo "=== admin ==="
find api/v1/admin -name "*_api.go" ! -name "*_test.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null
echo ""
echo "P2 System Modules:"
echo "=== system ==="
find api/v1/system -name "*_api.go" ! -name "*_test.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null
