#!/usr/bin/sh

#PRE COMMIT hook

git diff --cached --name-only | grep -E '*.yaml' |   GREP_COLOR='4;5;30;33' xargs grep --color --with-filename -n -i 'kind: Secret' && echo 'COMMIT REJECTED. Please encrypt before commiting' && exit 1

exit 0

