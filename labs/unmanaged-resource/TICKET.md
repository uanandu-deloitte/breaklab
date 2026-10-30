# INC-441: The MOTD update never actually applied

/infra/resources.yaml was updated with this week's message, and
infra-apply ran without any error. But /etc/motd.conf still shows last
week's text. Get it actually applying.

Get a shell: docker exec -it breaklab-unmanaged-resource bash
