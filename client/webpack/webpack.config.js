const path = require('path');
const webpack = require('webpack');
const ExtractTextPlugin = require('extract-text-webpack-plugin');

const rootPath = path.resolve(__dirname, "../../client");
const assetPath = path.resolve(__dirname, "../../assets");
const eslintPath = path.resolve(__dirname, '../../.eslintrc.json');

module.exports = [
  {
    entry: path.resolve(rootPath, 'index.jsx'),
    output: {
      path: path.resolve(assetPath, "js"),
      filename: 'bundle.js',
    },
    module: {
      loaders: [
        {
          test: /\.js[x]?$/,
          exclude: /node_modules/,
          loader: "babel-loader",
          query:{
            presets: ['react', ['es2015', { "modules": false }]],
          },
        },
        {
          enforce: 'pre',
          test: /\.js[x]?$/,
          exclude: /node_modules/,
          loader: "eslint-loader",
        },
      ],
    },
    plugins: [
      new webpack.LoaderOptionsPlugin({
        test: /\.js$/,
        options: {
          eslint: {
            configFile: eslintPath,
          },
        },
      }),
      new webpack.DefinePlugin({
        process: {
          env: {
            NODE_ENV: JSON.stringify(process.env.NODE_ENV),
          },
        },
      }),
    ],
    resolve: {
      extensions: ['.js', '.jsx'],
      modules: [rootPath, "node_modules"],
    },
  },
  {
    entry: path.resolve(rootPath, 'index.scss'),
    output: {
      path: path.resolve(assetPath, "css"),
      filename: 'main.css',
    },
    resolve: {
      modules: [rootPath, "node_modules"],
      extensions: [".css", ".scss"],
    },
    module: {
      rules: [
        {
          test: /\.css$/,
          loader: ExtractTextPlugin.extract({
            fallback: 'style-loader',
            use: "css-loader",
          }),
        },
        {
          test: /\.scss$/,
          loader: ExtractTextPlugin.extract({
            fallback: 'style-loader',
            use: ["css-loader", "sass-loader"],
          }),
        },
      ],
    },
    plugins: [
      new ExtractTextPlugin('main.css'),
    ],
    devtool: 'source-map',
  },
];
