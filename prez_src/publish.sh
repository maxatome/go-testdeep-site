#!/bin/sh -e

HOST=http://127.0.0.1:3999

TARGET="$(dirname $0)/../docs/prez"

curl -s $HOST/static/styles.css |
    perl -pE 'BEGIN { undef $/ } s,^ +/\* Add explicit links \*/.*?\n\n,,smg' > $TARGET/static/styles.css

for i in *.svg; do
    [ $i != colored-output.svg -a $i != logo.svg ] &&
        curl -s $HOST/$i > $TARGET/../images/$i
done

curl -s $HOST/static/slides.js | sed -E 's/\/(static\/)/\1/' > $TARGET/static/slides.js
curl -s $HOST/go-testdeep.slide |
    sed -E 's/\/(static\/)/\1/' |
    sed -E '/play\.js/d' |
    sed -E '/ +$/d' |
    sed -E 's/([-a-zA-Z0-9]+\.svg)/..\/images\/\1/' > $TARGET/index.html
