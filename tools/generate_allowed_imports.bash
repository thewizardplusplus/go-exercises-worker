#!/usr/bin/env bash

function wrap_lines_with_brackets() {
  declare -r left_bracket="$1"
  declare -r right_bracket="$2"

  sed --regexp-extended "s/(.*)/$left_bracket\1$right_bracket/"
}

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
  | wrap_lines_with_brackets '"' '"' \
  | paste --serial --delimiters ',' \
  | wrap_lines_with_brackets "[" "]"
