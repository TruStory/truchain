'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});
exports.customError = customError;
const errorEx = require('error-ex');

function customError(name, opts) {
    return errorEx(name, opts);
}

const errorTypes = exports.errorTypes = ['FileDoesNotExistError', 'NotAValidYAMLCustomTypeError', 'InvalidYAMLSuiteSchemaError', 'InvalidSuiteConfigurationError', 'DisabledSuiteError', 'NoSpecsFoundError', 'NoSpecFoundMatchingNameError', 'RequestBodyNotFoundError', 'ResponseStatusCodeDidNotMatchError', 'ResponseHeaderValueDidNotMatchError', 'YAMLSuiteLoadingError', 'SuiteConfigurationFailedError', 'BeforeAllHookError', 'AfterAllHookError', 'InvalidSpecificationSchemaError', 'BeforeEachHookError', 'AfterEachHookError', 'BeforeTestHookError', 'AfterTestHookError', 'InvalidRequestSpecificationError', 'InvalidRequestHeaderError', 'RequestBuilderError', 'RequestBodyBuilderError', 'JSONBodyParseError', 'LoadingSpecDependencySuiteError', 'ResponseJSONDataMismatchError', 'ResponseJSONSchemaValidationError', 'CustomResponseValidationError', 'LoopItemsBuilderError', 'SuiteCustomConfigurationError', 'CustomFunctionNotFoundInModuleError', 'ResponseCookieValueDidNotMatchError'];