# Test Allowed HTML/SPA Paths
GET {{host}}/
HTTP 200
[Asserts]
header "Cache-Control" contains "no-store, no-cache, must-revalidate"

GET {{host}}/login
HTTP 200
[Asserts]
header "Cache-Control" contains "no-store, no-cache, must-revalidate"

GET {{host}}/register
HTTP 200

GET {{host}}/leaderboard
HTTP 200

GET {{host}}/impressum
HTTP 200

GET {{host}}/privacy
HTTP 200

GET {{host}}/agb
HTTP 200

GET {{host}}/trade
HTTP 200

GET {{host}}/profile
HTTP 200

GET {{host}}/users/alex_123
HTTP 200

GET {{host}}/users/test-user-99
HTTP 200

# Test Allowed Static Files
GET {{host}}/robots.txt
HTTP 200

GET {{host}}/sitemap.xml
HTTP 200

GET {{host}}/icon.svg
HTTP 200
