import antfu from '@antfu/eslint-config'

export default antfu({
  unocss: true,
  vue: true,
  typescript: true,
}, {
  files: ['**/Api.ts'],
  rules: {
    'eslint-comments/no-unlimited-disable': ['off'],
  },
})
