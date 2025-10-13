module.exports = {
    extends: ['stylelint-config-standard-scss'],
    plugins: ['stylelint-scss'],
    rules: {
        'no-empty-source': null,
        'block-no-empty': true,
        'color-no-invalid-hex': true
    },
    overrides: [
        {
            files: ['**/*.scss'],
            customSyntax: 'postcss-scss'
        }
    ]
}
