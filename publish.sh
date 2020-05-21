#!/bin/sh -e

cd docs_src
hugo -d ../docs --minify

if git diff; then
    echo "No change detected after hugo generation"
else
    read -p "Do you want to commit these changes [y/N]: " ans
    case "$ans" in
        [yY]) git commit -a ;;
    esac
fi
