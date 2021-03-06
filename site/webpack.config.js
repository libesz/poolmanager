const autoprefixer = require('autoprefixer');
module.exports = [{
    entry: ['./app.scss', './app.js'],
    output: {
        filename: 'bundle.js',
    },
    "mode": "development",
    module: {
      rules: [
        {
          test: /\.scss$/,
          use: [
            {
              loader: 'file-loader',
              options: {
                name: 'bundle.css',
              },
            },
            { loader: 'extract-loader' },
            { loader: 'css-loader' },
            { loader: 'sass-loader',
              options: {
                // Prefer Dart Sass
                implementation: require('sass'),
                
                // See https://github.com/webpack-contrib/sass-loader/issues/804
                webpackImporter: false,
                sassOptions: {
                includePaths: ['./node_modules']
                },
              }
            },
            { loader: 'extract-loader' },
            { loader: 'css-loader' },
            { loader: 'postcss-loader',
              options: {
                plugins: () => [autoprefixer()]
              }
            },
            { loader: 'sass-loader',
              options: {
                sassOptions: {
                  includePaths: ['./node_modules']
                },
                // Prefer Dart Sass
                implementation: require('sass'),

                // See https://github.com/webpack-contrib/sass-loader/issues/804
                webpackImporter: false,
              }
            },
          ]
        },
        {
            test: /\.js$/,
            loader: 'babel-loader',
            query: {
              presets: ['@babel/preset-env'],
            },
        }
      ]
    },
  }];