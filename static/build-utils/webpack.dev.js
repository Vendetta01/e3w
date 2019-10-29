const Dotenv = require('dotenv-webpack');

module.exports = {
    mode: 'development',
    devtool: 'eval-source-map',
    output: {
	publicPath: '/'
    },
    plugins: [
	new Dotenv({
	    path: './.env.development',
	})
    ],
};

