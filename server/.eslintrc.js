module.exports = {
  env: {
    browser: true,
    commonjs: true,
    es2021: true,
    node: true,
    jest: true,
  },
  extends: 'airbnb-base',
  parserOptions: {
    ecmaVersion: 'latest',
  },
  rules: {
    camelcase: 'off',
    indent: 'off',
    'no-tabs': 'off',
    'no-underscore-dangle': 'off',
    'class-methods-use-this': 'off',
    'no-unused-vars': 'off',
    'prefer-destructuring': 'off',
    'linebreak-style': 'off',
  },
};
