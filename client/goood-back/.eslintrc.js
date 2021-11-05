module.exports = {
  root: true,
  env: {
    node: true
  },
  extends: [
    'plugin:vue/essential',
    '@vue/standard',
    '@vue/typescript/recommended'
  ],
  parserOptions: {
    ecmaVersion: 2020
  },
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    indent: 'off',
    quotes: 'off',
    'no-tabs': 'off',
    'padded-blocks': 'off',
    camelcase: 'off',
    'no-trailing-spaces': 'off',
    '@typescript-eslint/no-empty-function': 'off',
    'vue/script-indent': ['warn', 2, {
      baseIndent: 1
    }]
  }
}
