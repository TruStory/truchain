'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});

var _fs = require('fs');

var _fs2 = _interopRequireDefault(_fs);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _loader = require('./schema/yaml/loader');

var _errors = require('./errors');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

class SuiteDependency {

    constructor(filePath, parent) {
        this.file = filePath;
        this.parent = parent;
        this.targetConfiguration = {};
        this.commonHeaders = {};
    }

    loadFileAndValidateSchema() {
        this.parent.loadFileAndValidateSchema.call(this);
    }

    loadFile(file) {
        if (!_fs2.default.existsSync(file)) {
            let FileDoesNotExistError = (0, _errors.customError)('FileDoesNotExist');
            throw new FileDoesNotExistError(`Test suite file doesn't exist at '${file}'`);
        }

        try {
            return (0, _loader.loadYAML)(file, { encoding: 'UTF-8', customTypes: ['asyncFunction'] });
        } catch (err) {
            let YAMLSuiteLoadingError = (0, _errors.customError)('YAMLSuiteLoadingError');
            throw new YAMLSuiteLoadingError(`(${file}) \n ${err.message}`);
        }
    }

    resolveFile(filePath) {
        return this.parent.resolveFile.call(this, filePath);
    }

    async configure() {
        await this.parent.configure.call(this);
    }
}
exports.default = SuiteDependency;
module.exports = exports['default'];