# PROMAN
 

some simple web based UI to edit prometheuse.yml directly
  
## Pre-requisite
Make sure the prometheus use --web.enable-lifecycle flag, so this app can refresh the config file

  

## BUILDING

use the Makefile to build both the UI and backend.
To build dev SPA  
```
make build_web_dev
make build_api
```
To build production SPA
```
make build_web
make build_api
```
Make sure you have .env.production file properly made on ui directory

You will have dist directory containing the executable. You can run it using
```
cd dist 
./proman
```
## Database
by default this program uses mongodb as database for its user. The connection string is in cmd/api/config/api.json

You can create used dynamically by POST-ing REST to /auth/registeruser with payload 
```
{
"Username":"string",
"Password":"string"
}
```
