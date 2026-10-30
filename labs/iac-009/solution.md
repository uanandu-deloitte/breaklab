# Solution: The Config Says One Thing, The State Remembers Another

```bash
cd /infra
terraform plan
# ~ update in place
#   ~ content = "Welcome to Acme Checkout -- build 2026.09.14\n" -> "Welcome to Acme Checkout -- build 2026.10.20 (cart-fix)\n"

terraform apply -auto-approve
# writes the new banner text to /etc/motd-app/banner.txt and updates state

terraform plan -detailed-exitcode
# exit code 0 -- config, state, and the real file now all agree
```
