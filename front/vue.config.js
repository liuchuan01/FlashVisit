const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  lintOnSave: false,  // 禁用开发环境下的 ESLint 错误提示
  devServer: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
