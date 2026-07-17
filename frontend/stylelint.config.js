/** @type {import('stylelint').Config} */
export default {
  extends: ['stylelint-config-standard'],
  plugins: ['stylelint-order'],
  rules: {
    'selector-pseudo-class-no-unknown': true,
    'selector-max-id': 0,
    'order/properties-alphabetical-order': true,
    'block-no-empty': true,
    'selector-class-pattern': null,
  },
};
