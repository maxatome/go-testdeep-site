#!/bin/sh -e

set -o pipefail

cd "$(dirname "$0")"

HOST=http://127.0.0.1:3999

present &
trap "kill $!" EXIT INT
sleep 1

TARGET="$(dirname $0)/../docs/prez"
rm -rf "$TARGET"

mkdir -p "$TARGET/static"

curl -s "$HOST/static/styles.css" |
    perl -0777 -pE 's,^ +/\* Add explicit links \*/.*?\n\n,,smg' > "$TARGET/static/styles.css"

for i in *.svg; do
    [ "$i" != colored-output.svg -a "$i" != logo.svg ] &&
        curl -s "$HOST/$i" > "$TARGET/../images/$i"
done

curl -s "$HOST/static/slides.js" | sed -E 's,/(static/),\1,' > "$TARGET/static/slides.js"
curl -s "$HOST/go-testdeep.slide" |
    sed -E 's,/(static/),\1,' |
    sed -E '/play\.js/d' |
    sed -E '/ +$/d' |
    sed -E 's,(<head>),\1<link href=/images/favicon.png rel=icon type=image/png>,' |
    sed -E 's,([-a-zA-Z0-9]+\.svg),../images/\1,' > "$TARGET/index.html"
