const path = require('path');

module.exports = {
  mode: 'production', // or 'development'
  entry: './grpc/client.js',
  output: {
    filename: 'grpc-web-client.js',
    path: path.resolve(__dirname, 'public'),
  },
  resolve: {
    extensions: ['.ts', '.js'],
  },
  module: {
    rules: [
      {
        test: /\.ts$/,
        use: 'ts-loader', // Use ts-loader to transpile TypeScript
        exclude: /node_modules/,
      },
      {
        enforce: 'pre',
        test: /\.js$/,
        loader: 'source-map-loader', // Use source-map-loader for better debugging
      },
    ],
  },
  devtool: 'source-map', // Enable source maps for debugging
};
