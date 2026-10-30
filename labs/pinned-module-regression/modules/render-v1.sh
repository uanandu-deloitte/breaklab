render() {
  local msg="$1"
  # Bug: only keeps the first word of the message
  echo "${msg%% *}"
}
