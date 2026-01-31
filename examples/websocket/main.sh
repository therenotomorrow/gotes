#!/usr/bin/env bash

set -e

echo 'Receiving token via user registration'

token=$(curl 'http://127.0.0.1:8001/api/v1/users/register' \
  -X POST -H "Content-Type: application/json" \
  -d '{"name":"user","email":"user@memes.com","password":"supersecret"}' | jq -r '.user.token')

echo 'Send few events to server'

curl 'http://127.0.0.1:8001/api/v1/notes' \
  -X POST -H "Authorization: $token" -H 'Content-Type: application/json' \
  -d '{"title":"title 1","content":"content number 1"}' && echo

curl 'http://127.0.0.1:8001/api/v1/notes' \
  -X POST -H "Authorization: $token" -H 'Content-Type: application/json' \
  -d '{"title":"title 2","content":"content number 2"}' && echo

ident=$(curl 'http://127.0.0.1:8001/api/v1/notes' \
  -X POST -H "Authorization: $token" -H 'Content-Type: application/json' \
  -d '{"title":"title 3","content":"content number 3"}' | jq -r '.note.id.value')

curl "http://127.0.0.1:8001/api/v1/notes/$ident" \
  -X DELETE -H "Authorization: $token" && echo

wscat -c "ws://127.0.0.1:8001/api/v1/notes/events?token=$token"
#wscat -c "ws://127.0.0.01:8001/api/v1/notes/events?token=78f5e927-ef57-4b25-9a4c-11bc2b9d7858"
