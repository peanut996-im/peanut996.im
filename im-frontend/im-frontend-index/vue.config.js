// eslint-disable-next-line import/no-extraneous-dependencies
const webpack = require('webpack');
// const htmlPlugin = require('html-webpack-plugin');
const CompressionWebpackPlugin = require('compression-webpack-plugin');
// const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const os = require('os');

// eslint-disable-next-line no-unused-vars
const isDev = process.env.NODE_ENV === 'development';

// cdn链接
// eslint-disable-next-line no-unused-vars
const cdn = {
  css: [
    // antd css 由于引入失败只好放弃了antd的按需引入
  ],
  js: [
    // vue
    'https://cdn.bootcdn.net/ajax/libs/vue/2.6.10/vue.min.js',
    // vue-router
    'https://cdn.bootcdn.net/ajax/libs/vue-router/3.1.3/vue-router.min.js',
    // vuex
    'https://cdn.bootcdn.net/ajax/libs/vuex/3.1.2/vuex.min.js',
    // axios
    'https://cdn.bootcdn.net/ajax/libs/axios/0.18.0/axios.min.js',
    // moment
    'https://cdn.bootcdn.net/ajax/libs/moment.js/2.27.0/moment.min.js',
    // lodash
    'https://cdn.bootcdn.net/ajax/libs/lodash.js/4.17.20/lodash.min.js',
  ],
};

const htmlOptions = {
  favicon: 'src/assets/favicon.png',
  title: 'Peanut996.IM',
  filename: 'index.html',
  // vue default template
  template: 'public/index.html',
  inject: true,
};
// 配置cdn引入

module.exports = {
  /* eslint-disable */
    chainWebpack: (config) => {
        if (process.env.NODE_ENV === 'production' && process.env.VUE_APP_CDN === 'true') {
            const externals = {
                vue: 'Vue',
                axios: 'axios',
                'vue-router': 'VueRouter',
                vuex: 'Vuex',
                moment: 'moment',
                lodash: '_',
            };
            config.externals(externals);
            // 通过 html-webpack-plugin 将 cdn 注入到 index.html 之中
            config.plugin('html').tap((args) => {
                // eslint-disable-next-line no-param-reassign
                args[0].cdn = cdn;
                return args;
            });
        }
        config
            .plugin('html')
            .tap((args) => {
                args[0].template = htmlOptions.template;
                args[0].favicon = htmlOptions.favicon;
                args[0].title = htmlOptions.title;
                args[0].inject = htmlOptions.inject;
                return args;
            });
        if (process.env.NODE_ENV === 'production') {
            // increase build performance
            config
                .plugin('fork-ts-checker')
                .tap(args => {
                    let totalmem = Math.floor(os.totalmem() / 1024 / 1024); //get OS mem size
                    let allowUseMem = totalmem > 4096 ? 2048 : 1024;
                    // in vue-cli should args[0]['typescript'].memoryLimit
                    args[0].memoryLimit = allowUseMem;
                    args[0].workers = os.cpus().length;
                    return args
                });
        }
    },
    /* eslint-enable */
  configureWebpack: (config) => {
    // 代码 gzip
    const productionGzipExtensions = ['html', 'js', 'css'];
    // 开发模式下不走gzip
    if (!isDev) {
      config.plugins.push(
        new CompressionWebpackPlugin({
          filename: '[path].gz[query]',
          algorithm: 'gzip',
          test: new RegExp(`\\.(${productionGzipExtensions.join('|')})$`),
          threshold: 10240, // 只有大小大于该值的资源会被处理 10240
          minRatio: 0.8, // 只有压缩率小于这个值的资源才会被处理
          deleteOriginalAssets: false, // 删除原文件
        })
      );
    }
    // 不打包moment的语言包
    config.plugins.push(new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/));
    // 去除console
    if (process.env.NODE_ENV === 'production') {
      // eslint-disable-next-line no-param-reassign
      config.optimization.minimizer[0].options.terserOptions.compress.drop_console = true;
    }
  },
  css: {
    loaderOptions: {
      less: {
        lessOptions: {
          modifyVars: {
            'primary-color': '#09b955',
            // 'link-color': '#1DA57A',
            // 'border-radius-base': '2px',
          },
          javascriptEnabled: true,
        },
      },
      sass: {
        prependData: "@import '@/styles/index.scss';",
      },
    },
  },
  // webSocket本身不存在跨域问题，所以我们可以利用webSocket来进行非同源之间的通信。
  publicPath: './',
  productionSourceMap: false,
  devServer: {
    port: 8080,
    overlay: {
      warnings: false,
      errors: true,
    },
  },
  lintOnSave: true,
};
