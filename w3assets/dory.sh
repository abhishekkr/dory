
export LOCAL_AUTH_ADDR="http://127.0.0.1:8080"
export LOCAL_AUTH_TOKEN=""
export LOCAL_AUTH_PATH="pg"
export LOCAL_AUTH_TTL=60  ## default is 300 i.e. 5min
export LOCAL_AUTH_LS_HTTP_PARAMS='?keep=false'

seperator(){
  echo ""
  echo "--------------------------------------"
}

local-auth-unset(){
  unset LOCAL_AUTH_TOKEN
  echo $LOCAL_AUTH_TOKEN
}

local-auth-mount(){
  cat > /tmp/secret <<PEOF
{
  "username": "postgres",
  "password": "mostdiversedb"
}
PEOF

  export LOCAL_AUTH_TOKEN=$(curl -skL -X POST --data @/tmp/secret "${LOCAL_AUTH_ADDR}/local-auth/${LOCAL_AUTH_PATH}?ttlsecond=${LOCAL_AUTH_TTL}")

  echo "auth-token: "$LOCAL_AUTH_TOKEN

  seperator
}

local-auth-ls(){
  curl -skL \
      --header "X-DORY-TOKEN: ${LOCAL_AUTH_TOKEN}" \
      --request GET \
      "${LOCAL_AUTH_ADDR}/local-auth/${LOCAL_AUTH_PATH}${LOCAL_AUTH_LS_HTTP_PARAMS}"

  seperator
}

local-auth-unmount(){
  curl -skL \
      --header "X-DORY-TOKEN: ${LOCAL_AUTH_TOKEN}" \
      --request DELETE \
      ${LOCAL_AUTH_ADDR}/local-auth/${LOCAL_AUTH_PATH}

  seperator
}

local-auth-mount-persist(){
  cat > /tmp/secret <<PEOF
{
  "username": "postgres",
  "password": "mostdiversedb"
}
PEOF

  export LOCAL_AUTH_TOKEN=$(curl -skL -X POST --data @/tmp/secret "${LOCAL_AUTH_ADDR}/local-auth/${LOCAL_AUTH_PATH}?persist=true")

  echo "auth-token: "$LOCAL_AUTH_TOKEN

  seperator
}

local-auth-ls-persist(){
  curl -skL \
      --header "X-DORY-TOKEN: ${LOCAL_AUTH_TOKEN}" \
      --request GET \
      "${LOCAL_AUTH_ADDR}/local-auth/${LOCAL_AUTH_PATH}${LOCAL_AUTH_LS_HTTP_PARAMS}&persist=true"

  seperator
}

local-auth-unmount-persist(){
  curl -skL \
      --header "X-DORY-TOKEN: ${LOCAL_AUTH_TOKEN}" \
      --request DELETE \
      ${LOCAL_AUTH_ADDR}/local-auth/${LOCAL_AUTH_PATH}?persist=true

  seperator
}

