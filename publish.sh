#!/bin/sh -e

cd docs_src
hugo -d ../docs --minify

trap "rm -f /tmp/publish-$$.diff" EXIT INT

if git diff --exit-code --color=always > /tmp/publish-$$.diff; then
    echo "No change detected after hugo generation"
else
    less -R /tmp/publish-$$.diff
    read -p "Do you want to commit these changes [y/N]: " ans
    case "$ans" in
        [yY]) git commit -a ;;
    esac
fi
