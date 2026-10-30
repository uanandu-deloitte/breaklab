# INC-557: Orders stop completing, nothing is restarting

Support says new orders aren't completing, but `kubectl get pods` shows
nothing crashed or restarting. Find the actual failure and get orders
flowing again.

Get a shell: docker exec -it breaklab-backend-dependency-down bash
