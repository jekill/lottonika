module.exports = {
  devServer: {
    host: '0.0.0.0',
    port: '4444',
  },

  pluginOptions: {
    i18n: {
      locale: 'ru',
      fallbackLocale: 'en',
      localeDir: 'locales',
      enableInSFC: true
    }
  }
};
