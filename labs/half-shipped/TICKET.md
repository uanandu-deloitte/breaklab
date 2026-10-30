# INC-9000: The release that's only half true

Checkout still shows last week's message, and the release engineer
swears tonight's fix shipped clean — `infra-apply` ran without error,
and `/srv/release-repo`'s commit log looks fine at a glance too.
Something along the way didn't actually land. There may be more than
one thing wrong here.

Get a shell: docker exec -it breaklab-half-shipped bash
