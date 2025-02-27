/* eslint-disable */
// @generated by protobuf-ts 2.9.3 with parameter output_javascript,optimize_code_size,long_type_string,add_pb_suffix,ts_nocheck,eslint_disable
// @generated from protobuf file "runme/config/v1alpha1/config.proto" (package "runme.config.v1alpha1", syntax proto3)
// tslint:disable
// @ts-nocheck
/* eslint-disable */
// @generated by protobuf-ts 2.9.3 with parameter output_javascript,optimize_code_size,long_type_string,add_pb_suffix,ts_nocheck,eslint_disable
// @generated from protobuf file "runme/config/v1alpha1/config.proto" (package "runme.config.v1alpha1", syntax proto3)
// tslint:disable
// @ts-nocheck
import { MessageType } from "@protobuf-ts/runtime";
/**
 * @generated from protobuf enum runme.config.v1alpha1.Config.FilterType
 */
export var Config_FilterType;
(function (Config_FilterType) {
    /**
     * @generated from protobuf enum value: FILTER_TYPE_UNSPECIFIED = 0;
     */
    Config_FilterType[Config_FilterType["UNSPECIFIED"] = 0] = "UNSPECIFIED";
    /**
     * @generated from protobuf enum value: FILTER_TYPE_BLOCK = 1;
     */
    Config_FilterType[Config_FilterType["BLOCK"] = 1] = "BLOCK";
    /**
     * @generated from protobuf enum value: FILTER_TYPE_DOCUMENT = 2;
     */
    Config_FilterType[Config_FilterType["DOCUMENT"] = 2] = "DOCUMENT";
})(Config_FilterType || (Config_FilterType = {}));
// @generated message type with reflection information, may provide speed optimized methods
class Config$Type extends MessageType {
    constructor() {
        super("runme.config.v1alpha1.Config", [
            { no: 1, name: "project", kind: "message", oneof: "source", T: () => Config_Project },
            { no: 2, name: "filename", kind: "scalar", oneof: "source", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "env", kind: "message", T: () => Config_Env },
            { no: 5, name: "filters", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Config_Filter },
            { no: 7, name: "log", kind: "message", T: () => Config_Log }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message runme.config.v1alpha1.Config
 */
export const Config = new Config$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Config_Project$Type extends MessageType {
    constructor() {
        super("runme.config.v1alpha1.Config.Project", [
            { no: 1, name: "dir", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "find_repo_upward", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 3, name: "ignore_paths", kind: "scalar", jsonName: "ignore", repeat: 2 /*RepeatType.UNPACKED*/, T: 9 /*ScalarType.STRING*/ },
            { no: 4, name: "disable_gitignore", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message runme.config.v1alpha1.Config.Project
 */
export const Config_Project = new Config_Project$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Config_Filter$Type extends MessageType {
    constructor() {
        super("runme.config.v1alpha1.Config.Filter", [
            { no: 1, name: "type", kind: "enum", T: () => ["runme.config.v1alpha1.Config.FilterType", Config_FilterType, "FILTER_TYPE_"], options: { "buf.validate.field": { enum: { definedOnly: true } } } },
            { no: 2, name: "condition", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "1", maxLen: "1024" } } } }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message runme.config.v1alpha1.Config.Filter
 */
export const Config_Filter = new Config_Filter$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Config_Env$Type extends MessageType {
    constructor() {
        super("runme.config.v1alpha1.Config.Env", [
            { no: 1, name: "use_system_env", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "sources", kind: "scalar", repeat: 2 /*RepeatType.UNPACKED*/, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message runme.config.v1alpha1.Config.Env
 */
export const Config_Env = new Config_Env$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Config_Log$Type extends MessageType {
    constructor() {
        super("runme.config.v1alpha1.Config.Log", [
            { no: 1, name: "enabled", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "path", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "verbose", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message runme.config.v1alpha1.Config.Log
 */
export const Config_Log = new Config_Log$Type();
