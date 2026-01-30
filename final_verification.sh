#!/bin/bash
echo "=== Final Verification Report ==="
echo ""
echo "Checking for remaining shared.Success/SuccessData/Error calls in API files..."
echo ""

total=0
for dir in api/v1/*/; do
  count=$(find "$dir" -name "*_api.go" ! -name "*_test.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null | wc -l)
  if [ $count -gt 0 ]; then
    echo "$dir: $count files with shared calls"
    find "$dir" -name "*_api.go" ! -name "*_test.go" -exec grep -l "shared\.Success\|shared\.SuccessData\|shared\.Error" {} \; 2>/dev/null
    total=$((total + count))
  fi
done

echo ""
echo "Total files with shared calls: $total"
echo ""
echo "Checking compilation..."
go build ./api/v1/... 2>&1 | head -20
echo ""
echo "=== Verification Complete ==="
