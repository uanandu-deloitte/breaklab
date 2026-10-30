# INC-812: Maintenance mode won't turn off

We flipped `configwatcher` out of maintenance mode over an hour ago, but it's
still reporting maintenance mode internally. Make it reflect the current
configuration.

Get a shell: docker exec -it breaklab-stale-config bash
