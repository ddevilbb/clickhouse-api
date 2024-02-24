#!/usr/bin/env bash

log() {
  echo "START-SCRIPT: $1"
}

build() {
  log "Building server binary"
  go mod vendor
  go build -gcflags "all=-N -l" -mod vendor -o /go/src/app/build/cmd/api/main /go/src/app/cmd/api/main.go
  cp .env.dist ./build/cmd/api/.env
}

run() {
  log "Run server"

  log "Killing old server"
  killall dlv
  killall main

  log "Run in debug mode"
  /go/bin/dlv --continue --headless=true --listen=:2345 --api-version=2 --accept-multiclient exec /go/src/app/build/cmd/api/main &
}

rerun() {
  log "Rerun server"
  build
  run
}

hotReloading() {
  log "Run hotReloading"
  inotifywait --exclude "./vendor" -e "CREATE,MODIFY,DELETE,MOVED_TO,MOVED_FROM" -m -r ./ | (
    while true; do
      read path action file
      ext=${file: -3}
      if [[ "$ext" == ".go" ]]; then
        echo "$file"
      fi
    done
  ) | (
    WAITING=""
    while true; do
      file=""
      read -t 1 file
      if test -z "$file"; then
        if test ! -z "$WAITING"; then
          echo "CHANGED"
          WAITING=""
        fi
      else
        log "File ${file} changed" >>/tmp/filechanges.log
        WAITING=1
      fi
    done
  ) | (
    while true; do
      read TMP
      log "File Changed. Reloading..."
      rerun
    done
  )
}

initFileChangeLogger() {
  echo "" > /tmp/filechanges.log
  tail -f /tmp/filechanges.log &
}

initSqliteUsersDB() {
  mkdir -m 777 -p /sqlite/db
  sqlite3 /sqlite/db/users.db ".read /sqlite/users.sql"
}

main() {
  initFileChangeLogger
  build
  run
  initSqliteUsersDB
  hotReloading
}

main
