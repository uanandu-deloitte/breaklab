# INC-452: /infra/site.yaml was updated, nothing changed

Someone updated /infra/site.yaml with this quarter's message this
morning. infra-apply ran with no errors. /etc/site/deployed.conf still
shows an old message nobody can find in the file. Get the deployed
config to actually match /infra/site.yaml.

Get a shell: docker exec -it breaklab-infra-env-override bash
