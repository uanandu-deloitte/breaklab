# INC-812: Checkout failing, every Pod shows healthy

A security pass applied NetworkPolicy changes to this namespace this
morning. Since then, checkout has been intermittently failing even
though every Pod, Deployment, and Service looks completely healthy. Find
what's actually blocking traffic and get checkout working again.

Get a shell: docker exec -it breaklab-network-policy-block bash
