install-npm-deps:
	#from: https://github.com/material-components/material-components-web/blob/master/docs/getting-started.md
	npm i material-components-web
	npm install --save-dev webpack webpack-cli webpack-dev-server css-loader sass-loader sass extract-loader file-loader
	npm install --save-dev autoprefixer postcss-loader
	npm install --save-dev @babel/core babel-loader @babel/preset-env

build-site:
	npm run build
	cp dist/* ../pkg/webui/content/static/raw
	go generate github.com/libesz/poolmanager/pkg/webui/content/static