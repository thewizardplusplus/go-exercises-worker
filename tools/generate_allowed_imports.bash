#!/usr/bin/env bash

declare -ra disabledImports=(
  "builtin"
  "cmd"
  "internal"
  # TODO: complete the list of the disabled imports
)
go list "$(go env GOROOT)/src/..." \
  | grep \
    --invert-match \
    --perl-regexp "$(IFS="|"; echo "${disabledImports[*]}")" \
  | sed --regexp-extended 's/(.*)/"\1"/' \
  | paste --serial --delimiters ',' \
  | sed --regexp-extended 's/(.*)/[\1]/'
