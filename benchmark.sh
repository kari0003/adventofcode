#!/bin/zsh
if ! command -v hyperfine >/dev/null 2>&1; then
  echo "Please install hyperfine."
  exit 1
fi

execute () {
  sudo purge
  (
    trap time EXIT
    hyperfine \
      --warmup 5 \
      --min-runs 10 \
      --max-runs 300 \
      --command-name $1 \
      --prepare 'sync;' \
      'bun index.ts'
  )
  echo " "
}

execute bun

# if test -z $1; then
#   for dir in day_*; do
#    execute $dir
#   done
# else
#   execute $1
# fi