# INC-9200: The hotfix that didn't actually ship

Tonight's hotfix was supposedly merged, tagged, and built — but the
deployed artifact still shows the old broken checkout behavior. The
engineer on call swears the commit is in git and the pipeline ran
clean. There may be more than one thing wrong here.

Get a shell: docker exec -it breaklab-phantom-release bash
