#!/bin/sh

set -e

cd "$(dirname "$0")"
. ../util.sh
current_test=$(basename $(pwd))

export LOCAL_DIR=$(mktemp -d)

docker compose up -d --quiet-pull
sleep 5

docker compose exec backup backup --config-style=labels

sleep 5

expect_running_containers "2"

tmp_dir=$(mktemp -d)
tar -xvf "$LOCAL_DIR/test.tar.gz" -C "$tmp_dir"
if [ ! -f "$tmp_dir/backup/app_data/foo.txt" ]; then
  fail "Could not find expected file in untared archive."
fi

pass "Found relevant file in untared local backup."

# Backdate a dummy backup to verify pruning uses the rotation value from labels
touch -d "14 days ago" "$LOCAL_DIR/old.tar.gz"

docker compose exec backup backup --config-style=labels

if [ -f "$LOCAL_DIR/old.tar.gz" ]; then
  fail "Old backup has not been pruned."
fi

pass "Old backup has been pruned according to rotation label."
