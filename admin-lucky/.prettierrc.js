module.exports = {
    // Newline exceeding the maximum value
    printWidth: 148,
    // indent with 2 spaces
    tabWidth: 2,
    // do not use indentation, but use spaces
    useTabs: false,
    // There is no need for a semicolon at the end of the line
    semi: false,
    // use single quotes
    singleQuote: true,
    // Object keys are only quoted if necessary
    quoteProps: 'as-needed',
    // jsx does not use single quotes, but double quotes
    jsxSingleQuote: false,
    // Print trailing commas when possible on multiple lines. (For example, a single-line array will never end with a comma.) Optional value "<none|es5|all>", default none
    trailingComma: 'none',
    // put spaces between objects, array brackets and literals "{ foo: bar }"
    bracketSpacing: true,
    // The back angle brackets of the jsx tag need to wrap
    jsxBracketSameLine: true,
    bracketSameLine: true,
    // Arrow function, when there is always only one parameter, parentheses are also required, when the 'avoid' arrow function has only one parameter, the parentheses can be ignored
    arrowParens: 'always',
    // The scope of each file format is the entire content of the file
    rangeStart: 0,
    rangeEnd: Infinity,
    // No need to write @prettier at the beginning of the file
    requirePragma: false,
    // No need to automatically insert @prettier at the beginning of the file
    insertPragma: false,
    // use the default wrapping standard
    proseWrap: 'preserve',
    // Determine whether the html should break or not according to the display style
    htmlWhitespaceSensitivity: 'css',
    // Line breaks use lf ending with optional value "<auto|lf|crlf|cr>"
    endOfLine: 'auto'
}