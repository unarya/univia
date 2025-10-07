#!/bin/bash
set -e

FILES=(
  "roles"
  "users"
  "friends"
  "messages"
  "profiles"
  "permissions"
  "role_permissions"
  "access_tokens"
  "refesh_tokens"
  "posts"
  "comments"
  "commnent_likes"
  "categories"
  "post_categories"
  "post_likes"
  "media"
  "post_shares"
  "verification_codes"
  "notifications"
  "follows"
  "user_sessions"
)

i=1
for name in "${FILES[@]}"; do
  num=$(printf "%06d" $i)

  up_file=$(ls | grep -E "[0-9]+_create_${name}\.up\.sql$" || true)
  if [ -n "$up_file" ]; then
    new_up="${num}_create_${name}.up.sql"
    if [ "$up_file" != "$new_up" ]; then
      mv "$up_file" "$new_up"
      echo "Renamed: $up_file -> $new_up"
    fi
  fi

  down_file=$(ls | grep -E "[0-9]+_create_${name}\.down\.sql$" || true)
  if [ -n "$down_file" ]; then
    new_down="${num}_create_${name}.down.sql"
    if [ "$down_file" != "$new_down" ]; then
      mv "$down_file" "$new_down"
      echo "Renamed: $down_file -> $new_down"
    fi
  fi

  i=$((i+1))
done
