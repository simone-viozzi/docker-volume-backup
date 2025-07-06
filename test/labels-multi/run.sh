#!/bin/sh

set -e

cd "$(dirname "$0")"
. ../util.sh
current_test=$(basename $(pwd))

export LOCAL_DIR=$(mktemp -d)

docker compose up -d --quiet-pull
sleep 5

docker compose exec backup backup --config-style=labels

expect_running_containers "3"

tmp_dir=$(mktemp -d)
tar -xvf "$LOCAL_DIR/one/test.tar.gz" -C "$tmp_dir"
if [ ! -f "$tmp_dir/backup/vol_one/foo.txt" ]; then
  fail "Could not find file from volume one."
fi

tmp_dir2=$(mktemp -d)
tar -xvf "$LOCAL_DIR/two/test.tar.gz" -C "$tmp_dir2"
if [ ! -f "$tmp_dir2/backup/vol_two/bar.txt" ]; then
  fail "Could not find file from volume two."
fi

pass "Backups for all labeled volumes were created."
