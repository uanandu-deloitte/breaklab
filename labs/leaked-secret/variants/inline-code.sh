#!/usr/bin/env bash
set -euo pipefail
cd /srv/repo
printf '#!/usr/bin/env bash\nAPI_KEY=sk-FAKE9876543210\ncurl -H "Authorization: Bearer $API_KEY" https://example.invalid\n' > deploy.sh
git add -A && git commit -q -m "Add deploy script"
printf '#!/usr/bin/env bash\nAPI_KEY="$(cat /run/secrets/api_key)"\ncurl -H "Authorization: Bearer $API_KEY" https://example.invalid\n' > deploy.sh
git add -A && git commit -q -m "Load API key from a secrets file instead of hardcoding it"
