const path = require('path');

module.exports = {
    entry: './src/browser_shim.ts',
    node: {
      fs: 'empty'
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: [{
                    loader: 'ts-loader',
                }],
                exclude: /node_modules/
            },
        ],
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
    },
    output: {
        filename: 'browser_shim.js',
        path: path.resolve(__dirname, 'dist'),
    }
};
