import { SignedOrder } from '@0x/order-utils';
import { AcceptedOrderInfo, Config, ContractEvent, GetOrdersResponse, OrderEvent, OrderInfo, RejectedOrderInfo, Stats, ValidationResults, WrapperAcceptedOrderInfo, WrapperConfig, WrapperContractEvent, WrapperGetOrdersResponse, WrapperOrderEvent, WrapperOrderInfo, WrapperRejectedOrderInfo, WrapperSignedOrder, WrapperStats, WrapperValidationResults } from './types';
export declare function configToWrapperConfig(config: Config): WrapperConfig;
export declare function wrapperSignedOrderToSignedOrder(wrapperSignedOrder: WrapperSignedOrder): SignedOrder;
export declare function wrapperContractEventsToContractEvents(wrapperContractEvents: WrapperContractEvent[]): ContractEvent[];
export declare function signedOrderToWrapperSignedOrder(signedOrder: SignedOrder): WrapperSignedOrder;
export declare function wrapperOrderEventToOrderEvent(wrapperOrderEvent: WrapperOrderEvent): OrderEvent;
export declare function orderEventsHandlerToWrapperOrderEventsHandler(orderEventsHandler: (events: OrderEvent[]) => void): (events: WrapperOrderEvent[]) => void;
export declare function wrapperValidationResultsToValidationResults(wrapperValidationResults: WrapperValidationResults): ValidationResults;
export declare function wrapperAcceptedOrderInfoToAcceptedOrderInfo(wrapperAcceptedOrderInfo: WrapperAcceptedOrderInfo): AcceptedOrderInfo;
export declare function wrapperRejectedOrderInfoToRejectedOrderInfo(wrapperRejectedOrderInfo: WrapperRejectedOrderInfo): RejectedOrderInfo;
export declare function wrapperStatsToStats(wrapperStats: WrapperStats): Stats;
export declare function wrapperGetOrdersResponseToGetOrdersResponse(wrapperGetOrdersResponse: WrapperGetOrdersResponse): GetOrdersResponse;
export declare function wrapperOrderInfoToOrderInfo(wrapperOrderInfo: WrapperOrderInfo): OrderInfo;
//# sourceMappingURL=wrapper_conversion.d.ts.map