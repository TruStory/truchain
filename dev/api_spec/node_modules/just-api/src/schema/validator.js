'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});
exports.validateJSONSchema = validateJSONSchema;

var _jsonschema = require('jsonschema');

var _fs = require('fs');

var _fs2 = _interopRequireDefault(_fs);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function validateJSONSchema(dataObject, schemaFilePath) {
    const expectedSchema = JSON.parse(_fs2.default.readFileSync(schemaFilePath, 'utf8'));

    return (0, _jsonschema.validate)(dataObject, expectedSchema);
}