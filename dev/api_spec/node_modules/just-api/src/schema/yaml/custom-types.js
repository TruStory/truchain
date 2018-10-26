'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});
exports.asyncFunction = asyncFunction;

var _jsYaml = require('js-yaml');

var _jsYaml2 = _interopRequireDefault(_jsYaml);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

const esprima = require('esprima');


const AsyncFunction = Object.getPrototypeOf(async function () {}).constructor;

function asyncFunction() {
    return new _jsYaml2.default.Type('!js/asyncFunction', {
        kind: 'scalar',
        construct: function (data) {
            let functionBody = data || '';
            return constructJavascriptFunction(functionBody);
        },
        instanceOf: AsyncFunction
    });
}

function constructJavascriptFunction(data) {
    let source = '(' + data + ')',
        ast = esprima.parse(source, { range: true }),
        params = [],
        body;

    if (ast.type !== 'Program' || ast.body.length !== 1 || ast.body[0].type !== 'ExpressionStatement' || ast.body[0].expression.type !== 'ArrowFunctionExpression' && ast.body[0].expression.type !== 'FunctionExpression') {
        throw new Error(`Failed to resolve async function with body ${data}`);
    }

    ast.body[0].expression.params.forEach(function (param) {
        params.push(param.name);
    });

    body = ast.body[0].expression.body.range;

    return new AsyncFunction(params, source.slice(body[0] + 1, body[1] - 1));
}