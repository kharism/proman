clean:
	rm -r dist
	rm -r ui/dist
build_web_dev:
	mkdir -p dist
	cd ui && yarn install && yarn build --mode=development
	cp -r ui/dist/* dist/
build_web:
	mkdir -p dist
	cd ui && yarn install && yarn build
	cp -r ui/dist/* dist/
build_api:
	mkdir -p dist
	cd cmd/api && go build -o proman main.go
	cp cmd/api/proman ./dist
	cp -r cmd/api/config ./dist