"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var OrderEventKind;
(function (OrderEventKind) {
    OrderEventKind["Invalid"] = "INVALID";
    OrderEventKind["Added"] = "ADDED";
    OrderEventKind["Filled"] = "FILLED";
    OrderEventKind["FullyFilled"] = "FULLY_FILLED";
    OrderEventKind["Cancelled"] = "CANCELLED";
    OrderEventKind["Expired"] = "EXPIRED";
    OrderEventKind["Unfunded"] = "UNFUNDED";
    OrderEventKind["FillabilityIncreased"] = "FILLABILITY_INCREASED";
})(OrderEventKind = exports.OrderEventKind || (exports.OrderEventKind = {}));
var RejectedKind;
(function (RejectedKind) {
    RejectedKind["ZeroexValidation"] = "ZEROEX_VALIDATION";
    RejectedKind["MeshError"] = "MESH_ERROR";
    RejectedKind["MeshValidation"] = "MESH_VALIDATION";
})(RejectedKind = exports.RejectedKind || (exports.RejectedKind = {}));
var RejectedCode;
(function (RejectedCode) {
    RejectedCode["InternalError"] = "InternalError";
    RejectedCode["MaxOrderSizeExceeded"] = "MaxOrderSizeExceeded";
    RejectedCode["OrderAlreadyStored"] = "OrderAlreadyStored";
    RejectedCode["OrderForIncorrectNetwork"] = "OrderForIncorrectNetwork";
    RejectedCode["NetworkRequestFailed"] = "NetworkRequestFailed";
    RejectedCode["OrderHasInvalidMakerAssetAmount"] = "OrderHasInvalidMakerAssetAmount";
    RejectedCode["OrderHasInvalidTakerAssetAmount"] = "OrderHasInvalidTakerAssetAmount";
    RejectedCode["OrderExpired"] = "OrderExpired";
    RejectedCode["OrderFullyFilled"] = "OrderFullyFilled";
    RejectedCode["OrderCancelled"] = "OrderCancelled";
    RejectedCode["OrderUnfunded"] = "OrderUnfunded";
    RejectedCode["OrderHasInvalidMakerAssetData"] = "OrderHasInvalidMakerAssetData";
    RejectedCode["OrderHasInvalidTakerAssetData"] = "OrderHasInvalidTakerAssetData";
    RejectedCode["OrderHasInvalidSignature"] = "OrderHasInvalidSignature";
})(RejectedCode = exports.RejectedCode || (exports.RejectedCode = {}));
