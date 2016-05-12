#! /bin/bash

checks=$(guru -scope  . implements "check.go:#332" | sed -n -e 's/.* //; s/*//; 2,$p' | grep -v -E "^[^*A-Z]")

{
echo '<!DOCTYPE html>'
echo '<html><head><title>Availbale Checks</title><meta charset="UTF-8"></head>'
echo "<body><h1>Available Checks</h1>"

echo "<p>ht version $(git describe)</p>"

for c in $checks Condition; do
    [[ -n "$c" ]] || continue
    if [[ "$c" = "Condition" ]];then
	echo "<hr>"
	echo "<p>Type Condition is not a Check but it is used so often"
	echo "   in checks that it is worth describing here.</p>"
    fi
    echo "<h2>$c</h2>"
    echo "<pre>"
    go doc "$c" \
	| sed -e "s/&/\\&amp;/g; s/</\\&lt;/g; s/>/\\&gt;/g"
    echo "</pre>"
done | grep -v -E "^func " | sed -e "s/\t/        /; s/^    //; s/^}$/}\n/"

echo "</body></html>"
} > Checks.html

wkhtmltopdf --minimum-font-size 16 -s A4 -L 24 -R 24 -T 24 -B 24 Checks.html Checks.pdf
