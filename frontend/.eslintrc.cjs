const config = {
    extends: [
        'airbnb',
        'airbnb/hooks',
        'airbnb-typescript',
        'plugin:prettier/recommended',
        'prettier',
    ],
    parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.web.json'],
    },
    plugins: ['perfectionist'],
    rules: {
        'perfectionist/sort-imports': [
            'error',
            {
                customGroups: { type: {}, value: {} },
                environment: 'node',
                fallbackSort: { type: 'unsorted' },
                groups: [
                    'type',
                    ['builtin', 'external'],
                    'internal-type',
                    'internal',
                    ['parent-type', 'sibling-type', 'index-type'],
                    ['parent', 'sibling', 'index'],
                    'object',
                    'unknown',
                ],
                ignoreCase: true,
                internalPattern: ['^~/.+'],
                maxLineLength: undefined,
                newlinesBetween: 'always',
                order: 'asc',
                partitionByComment: false,
                partitionByNewLine: false,
                specialCharacters: 'keep',
                type: 'alphabetical',
            },
        ],
        'perfectionist/sort-objects': 'off',
        'react/react-in-jsx-scope': 'off',
        'import/no-extraneous-dependencies': 'warn',
        'import/prefer-default-export': 'off',
        'no-param-reassign': 'off',
        'import/order': 'off',
        'react/jsx-props-no-spreading': 'off',
    },
};

module.exports = config;
