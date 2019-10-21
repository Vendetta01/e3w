var path = require('path');
var HtmlWebpackPlugin = require('html-webpack-plugin')
var webpack = require('webpack')

var plugins = [
    new HtmlWebpackPlugin({
        filename: 'index.html',
        template: './src/index.html',
        inject: false,
    }),
    new HtmlWebpackPlugin({
        filename: 'login.html',
        template: './src/login.html',
        inject: false,
    })
]

process.env.NODE_ENV === 'production' ? plugins.push(new webpack.DefinePlugin({
    "process.env": {
        NODE_ENV: JSON.stringify("production")
    }
})) : null

module.exports = {
    devtool: process.env.NODE_ENV === 'production' ? '' : 'inline-source-map',
    entry: './src/entry.jsx',
    output: {
        path: path.join(__dirname, '/dist'),
        filename: 'bundle.js'
    },
    resolve: {
	extensions: ['.js', '.jsx'],
	enforceExtension: false
    },
    module: {
        rules: [{
            test: /.jsx$/,
	    loader: 'babel-loader',
            exclude: /node_modules/,
            query: {
                presets: ['@babel/preset-react', '@babel/preset-env'],
                plugins: [['import', { "libraryName": "antd"}], '@babel/plugin-proposal-class-properties']
            }
        }, {
            test: /\.css$/,
	    use : [
		'style-loader',
		'css-loader'
	    ]
        }]
    },
    plugins: plugins
}
