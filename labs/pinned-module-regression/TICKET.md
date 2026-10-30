# INC-467: The deployed message is wrong, the source isn't

/infra/site.yaml declares the correct order-confirmation message, and
infra-apply ran with no errors. But /etc/site/deployed.conf doesn't
match what the source says. Get the deployed config to actually match
/infra/site.yaml.

Get a shell: docker exec -it breaklab-pinned-module-regression bash
