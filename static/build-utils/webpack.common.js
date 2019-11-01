const path = require('path');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');


module.exports = {
    entry: './src/entry.jsx',
    output: {
	filename: 'bundle.js',
	path: path.resolve(__dirname, '../', 'dist')
    },
    resolve: {
	extensions: ['.js', '.jsx'],
	enforceExtension: false
    },
    module: {
	rules: [
	    {
		test: /\.(js|jsx)$/,
		exclude: /node_modules/,
		use: {
		    loader: 'babel-loader',
		    options: {
			presets: ['@babel/preset-env', '@babel/preset-react'],
			plugins: [
			    ['import', {
				"libraryName": "antd",
				"libraryDirectory": "es",
				"style": "true"
			    }, "ant"],
			    '@babel/plugin-proposal-class-properties'
			]
		    }
		}
	    },
	    {
		test: /\.css$/,
		use: [
		    'style-loader',
		    {
			loader: MiniCssExtractPlugin.loader
		    },
		    'css-loader'
		]
	    }
	]
    },
    plugins: [
	new CleanWebpackPlugin(),
	new HtmlWebpackPlugin({
	    filename: 'index.html',
	    template: './src/index.html'
	}),
	new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/),
	new MiniCssExtractPlugin({
	    filename: 'style.css',
	})
    ],
    devServer: {
	contentBase: './dist'
    }
};

