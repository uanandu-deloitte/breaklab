# render-v2.sh — the current rendering logic. Adds an explicit
# environment label that v1 never produced.
render() {
  local seed="$1"
  echo "Release ${seed} [rendered by v2]"
}
