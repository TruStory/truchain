'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});
exports.loadYAML = loadYAML;

var _jsYaml = require('js-yaml');

var _jsYaml2 = _interopRequireDefault(_jsYaml);

var _fs = require('fs');

var _fs2 = _interopRequireDefault(_fs);

var _customTypes = require('./custom-types');

var _errors = require('./../../errors');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function loadYAML(filePath, options) {
    let customSchema;

    if (options.customTypes) {
        let typeDefinitions = [];

        for (let type of options.customTypes) {
            typeDefinitions.push(getCustomYAMLSchema(type));
        }

        customSchema = _jsYaml2.default.Schema.create(typeDefinitions);
    }

    try {
        return _jsYaml2.default.load(_fs2.default.readFileSync(filePath, { encoding: options.encoding }), { schema: customSchema });
    } catch (err) {
        throw err;
    }
}

function getCustomYAMLSchema(customType) {
    let yamlType;

    switch (customType) {
        case 'asyncFunction':
            yamlType = (0, _customTypes.asyncFunction)();
            break;
        default:
            throw new _errors.NotAValidYAMLCustomType(`${customType} is not a valid custom YAML type`);
    }

    return yamlType;
}