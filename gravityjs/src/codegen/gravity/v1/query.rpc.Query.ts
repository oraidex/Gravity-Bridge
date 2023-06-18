import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryCurrentValsetRequest, QueryCurrentValsetResponse, QueryValsetRequestRequest, QueryValsetRequestResponse, QueryValsetConfirmRequest, QueryValsetConfirmResponse, QueryValsetConfirmsByNonceRequest, QueryValsetConfirmsByNonceResponse, QueryLastValsetRequestsRequest, QueryLastValsetRequestsResponse, QueryLastPendingValsetRequestByAddrRequest, QueryLastPendingValsetRequestByAddrResponse, QueryLastPendingBatchRequestByAddrRequest, QueryLastPendingBatchRequestByAddrResponse, QueryLastPendingLogicCallByAddrRequest, QueryLastPendingLogicCallByAddrResponse, QueryLastEventNonceByAddrRequest, QueryLastEventNonceByAddrResponse, QueryLastERC721EventNonceByAddrRequest, QueryLastERC721EventNonceByAddrResponse, QueryBatchFeeRequest, QueryBatchFeeResponse, QueryOutgoingTxBatchesRequest, QueryOutgoingTxBatchesResponse, QueryOutgoingLogicCallsRequest, QueryOutgoingLogicCallsResponse, QueryBatchRequestByNonceRequest, QueryBatchRequestByNonceResponse, QueryBatchConfirmsRequest, QueryBatchConfirmsResponse, QueryLogicConfirmsRequest, QueryLogicConfirmsResponse, QueryERC20ToDenomRequest, QueryERC20ToDenomResponse, QueryDenomToERC20Request, QueryDenomToERC20Response, QueryLastObservedEthBlockRequest, QueryLastObservedEthBlockResponse, QueryLastObservedEthNonceRequest, QueryLastObservedEthNonceResponse, QueryLastObservedERC721EthNonceRequest, QueryLastObservedERC721EthNonceResponse, QueryAttestationsRequest, QueryAttestationsResponse, QueryERC721AttestationsRequest, QueryERC721AttestationsResponse, QueryDelegateKeysByValidatorAddress, QueryDelegateKeysByValidatorAddressResponse, QueryDelegateKeysByEthAddress, QueryDelegateKeysByEthAddressResponse, QueryDelegateKeysByOrchestratorAddress, QueryDelegateKeysByOrchestratorAddressResponse, QueryPendingSendToEth, QueryPendingSendToEthResponse, QueryPendingIbcAutoForwards, QueryPendingIbcAutoForwardsResponse, QueryPendingERC721IbcAutoForwardsRequest, QueryPendingERC721IbcAutoForwardsResponse } from "./query";
/** Query defines the gRPC querier service */
export interface Query {
  /** Deployments queries deployments */
  params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
  currentValset(request?: QueryCurrentValsetRequest): Promise<QueryCurrentValsetResponse>;
  valsetRequest(request: QueryValsetRequestRequest): Promise<QueryValsetRequestResponse>;
  valsetConfirm(request: QueryValsetConfirmRequest): Promise<QueryValsetConfirmResponse>;
  valsetConfirmsByNonce(request: QueryValsetConfirmsByNonceRequest): Promise<QueryValsetConfirmsByNonceResponse>;
  lastValsetRequests(request?: QueryLastValsetRequestsRequest): Promise<QueryLastValsetRequestsResponse>;
  lastPendingValsetRequestByAddr(request: QueryLastPendingValsetRequestByAddrRequest): Promise<QueryLastPendingValsetRequestByAddrResponse>;
  lastPendingBatchRequestByAddr(request: QueryLastPendingBatchRequestByAddrRequest): Promise<QueryLastPendingBatchRequestByAddrResponse>;
  lastPendingLogicCallByAddr(request: QueryLastPendingLogicCallByAddrRequest): Promise<QueryLastPendingLogicCallByAddrResponse>;
  lastEventNonceByAddr(request: QueryLastEventNonceByAddrRequest): Promise<QueryLastEventNonceByAddrResponse>;
  lastERC721EventNonceByAddr(request: QueryLastERC721EventNonceByAddrRequest): Promise<QueryLastERC721EventNonceByAddrResponse>;
  batchFees(request?: QueryBatchFeeRequest): Promise<QueryBatchFeeResponse>;
  outgoingTxBatches(request?: QueryOutgoingTxBatchesRequest): Promise<QueryOutgoingTxBatchesResponse>;
  outgoingLogicCalls(request?: QueryOutgoingLogicCallsRequest): Promise<QueryOutgoingLogicCallsResponse>;
  batchRequestByNonce(request: QueryBatchRequestByNonceRequest): Promise<QueryBatchRequestByNonceResponse>;
  batchConfirms(request: QueryBatchConfirmsRequest): Promise<QueryBatchConfirmsResponse>;
  logicConfirms(request: QueryLogicConfirmsRequest): Promise<QueryLogicConfirmsResponse>;
  eRC20ToDenom(request: QueryERC20ToDenomRequest): Promise<QueryERC20ToDenomResponse>;
  denomToERC20(request: QueryDenomToERC20Request): Promise<QueryDenomToERC20Response>;
  getLastObservedEthBlock(request: QueryLastObservedEthBlockRequest): Promise<QueryLastObservedEthBlockResponse>;
  getLastObservedEthNonce(request: QueryLastObservedEthNonceRequest): Promise<QueryLastObservedEthNonceResponse>;
  getLastObservedERC721EthNonce(request?: QueryLastObservedERC721EthNonceRequest): Promise<QueryLastObservedERC721EthNonceResponse>;
  getAttestations(request: QueryAttestationsRequest): Promise<QueryAttestationsResponse>;
  getERC721Attestations(request: QueryERC721AttestationsRequest): Promise<QueryERC721AttestationsResponse>;
  getDelegateKeyByValidator(request: QueryDelegateKeysByValidatorAddress): Promise<QueryDelegateKeysByValidatorAddressResponse>;
  getDelegateKeyByEth(request: QueryDelegateKeysByEthAddress): Promise<QueryDelegateKeysByEthAddressResponse>;
  getDelegateKeyByOrchestrator(request: QueryDelegateKeysByOrchestratorAddress): Promise<QueryDelegateKeysByOrchestratorAddressResponse>;
  getPendingSendToEth(request: QueryPendingSendToEth): Promise<QueryPendingSendToEthResponse>;
  getPendingIbcAutoForwards(request: QueryPendingIbcAutoForwards): Promise<QueryPendingIbcAutoForwardsResponse>;
  getPendingERC721IbcAutoForwards(request: QueryPendingERC721IbcAutoForwardsRequest): Promise<QueryPendingERC721IbcAutoForwardsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.params = this.params.bind(this);
    this.currentValset = this.currentValset.bind(this);
    this.valsetRequest = this.valsetRequest.bind(this);
    this.valsetConfirm = this.valsetConfirm.bind(this);
    this.valsetConfirmsByNonce = this.valsetConfirmsByNonce.bind(this);
    this.lastValsetRequests = this.lastValsetRequests.bind(this);
    this.lastPendingValsetRequestByAddr = this.lastPendingValsetRequestByAddr.bind(this);
    this.lastPendingBatchRequestByAddr = this.lastPendingBatchRequestByAddr.bind(this);
    this.lastPendingLogicCallByAddr = this.lastPendingLogicCallByAddr.bind(this);
    this.lastEventNonceByAddr = this.lastEventNonceByAddr.bind(this);
    this.lastERC721EventNonceByAddr = this.lastERC721EventNonceByAddr.bind(this);
    this.batchFees = this.batchFees.bind(this);
    this.outgoingTxBatches = this.outgoingTxBatches.bind(this);
    this.outgoingLogicCalls = this.outgoingLogicCalls.bind(this);
    this.batchRequestByNonce = this.batchRequestByNonce.bind(this);
    this.batchConfirms = this.batchConfirms.bind(this);
    this.logicConfirms = this.logicConfirms.bind(this);
    this.eRC20ToDenom = this.eRC20ToDenom.bind(this);
    this.denomToERC20 = this.denomToERC20.bind(this);
    this.getLastObservedEthBlock = this.getLastObservedEthBlock.bind(this);
    this.getLastObservedEthNonce = this.getLastObservedEthNonce.bind(this);
    this.getLastObservedERC721EthNonce = this.getLastObservedERC721EthNonce.bind(this);
    this.getAttestations = this.getAttestations.bind(this);
    this.getERC721Attestations = this.getERC721Attestations.bind(this);
    this.getDelegateKeyByValidator = this.getDelegateKeyByValidator.bind(this);
    this.getDelegateKeyByEth = this.getDelegateKeyByEth.bind(this);
    this.getDelegateKeyByOrchestrator = this.getDelegateKeyByOrchestrator.bind(this);
    this.getPendingSendToEth = this.getPendingSendToEth.bind(this);
    this.getPendingIbcAutoForwards = this.getPendingIbcAutoForwards.bind(this);
    this.getPendingERC721IbcAutoForwards = this.getPendingERC721IbcAutoForwards.bind(this);
  }
  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new _m0.Reader(data)));
  }
  currentValset(request: QueryCurrentValsetRequest = {}): Promise<QueryCurrentValsetResponse> {
    const data = QueryCurrentValsetRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "CurrentValset", data);
    return promise.then(data => QueryCurrentValsetResponse.decode(new _m0.Reader(data)));
  }
  valsetRequest(request: QueryValsetRequestRequest): Promise<QueryValsetRequestResponse> {
    const data = QueryValsetRequestRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "ValsetRequest", data);
    return promise.then(data => QueryValsetRequestResponse.decode(new _m0.Reader(data)));
  }
  valsetConfirm(request: QueryValsetConfirmRequest): Promise<QueryValsetConfirmResponse> {
    const data = QueryValsetConfirmRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "ValsetConfirm", data);
    return promise.then(data => QueryValsetConfirmResponse.decode(new _m0.Reader(data)));
  }
  valsetConfirmsByNonce(request: QueryValsetConfirmsByNonceRequest): Promise<QueryValsetConfirmsByNonceResponse> {
    const data = QueryValsetConfirmsByNonceRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "ValsetConfirmsByNonce", data);
    return promise.then(data => QueryValsetConfirmsByNonceResponse.decode(new _m0.Reader(data)));
  }
  lastValsetRequests(request: QueryLastValsetRequestsRequest = {}): Promise<QueryLastValsetRequestsResponse> {
    const data = QueryLastValsetRequestsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "LastValsetRequests", data);
    return promise.then(data => QueryLastValsetRequestsResponse.decode(new _m0.Reader(data)));
  }
  lastPendingValsetRequestByAddr(request: QueryLastPendingValsetRequestByAddrRequest): Promise<QueryLastPendingValsetRequestByAddrResponse> {
    const data = QueryLastPendingValsetRequestByAddrRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "LastPendingValsetRequestByAddr", data);
    return promise.then(data => QueryLastPendingValsetRequestByAddrResponse.decode(new _m0.Reader(data)));
  }
  lastPendingBatchRequestByAddr(request: QueryLastPendingBatchRequestByAddrRequest): Promise<QueryLastPendingBatchRequestByAddrResponse> {
    const data = QueryLastPendingBatchRequestByAddrRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "LastPendingBatchRequestByAddr", data);
    return promise.then(data => QueryLastPendingBatchRequestByAddrResponse.decode(new _m0.Reader(data)));
  }
  lastPendingLogicCallByAddr(request: QueryLastPendingLogicCallByAddrRequest): Promise<QueryLastPendingLogicCallByAddrResponse> {
    const data = QueryLastPendingLogicCallByAddrRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "LastPendingLogicCallByAddr", data);
    return promise.then(data => QueryLastPendingLogicCallByAddrResponse.decode(new _m0.Reader(data)));
  }
  lastEventNonceByAddr(request: QueryLastEventNonceByAddrRequest): Promise<QueryLastEventNonceByAddrResponse> {
    const data = QueryLastEventNonceByAddrRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "LastEventNonceByAddr", data);
    return promise.then(data => QueryLastEventNonceByAddrResponse.decode(new _m0.Reader(data)));
  }
  lastERC721EventNonceByAddr(request: QueryLastERC721EventNonceByAddrRequest): Promise<QueryLastERC721EventNonceByAddrResponse> {
    const data = QueryLastERC721EventNonceByAddrRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "LastERC721EventNonceByAddr", data);
    return promise.then(data => QueryLastERC721EventNonceByAddrResponse.decode(new _m0.Reader(data)));
  }
  batchFees(request: QueryBatchFeeRequest = {}): Promise<QueryBatchFeeResponse> {
    const data = QueryBatchFeeRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "BatchFees", data);
    return promise.then(data => QueryBatchFeeResponse.decode(new _m0.Reader(data)));
  }
  outgoingTxBatches(request: QueryOutgoingTxBatchesRequest = {}): Promise<QueryOutgoingTxBatchesResponse> {
    const data = QueryOutgoingTxBatchesRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "OutgoingTxBatches", data);
    return promise.then(data => QueryOutgoingTxBatchesResponse.decode(new _m0.Reader(data)));
  }
  outgoingLogicCalls(request: QueryOutgoingLogicCallsRequest = {}): Promise<QueryOutgoingLogicCallsResponse> {
    const data = QueryOutgoingLogicCallsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "OutgoingLogicCalls", data);
    return promise.then(data => QueryOutgoingLogicCallsResponse.decode(new _m0.Reader(data)));
  }
  batchRequestByNonce(request: QueryBatchRequestByNonceRequest): Promise<QueryBatchRequestByNonceResponse> {
    const data = QueryBatchRequestByNonceRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "BatchRequestByNonce", data);
    return promise.then(data => QueryBatchRequestByNonceResponse.decode(new _m0.Reader(data)));
  }
  batchConfirms(request: QueryBatchConfirmsRequest): Promise<QueryBatchConfirmsResponse> {
    const data = QueryBatchConfirmsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "BatchConfirms", data);
    return promise.then(data => QueryBatchConfirmsResponse.decode(new _m0.Reader(data)));
  }
  logicConfirms(request: QueryLogicConfirmsRequest): Promise<QueryLogicConfirmsResponse> {
    const data = QueryLogicConfirmsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "LogicConfirms", data);
    return promise.then(data => QueryLogicConfirmsResponse.decode(new _m0.Reader(data)));
  }
  eRC20ToDenom(request: QueryERC20ToDenomRequest): Promise<QueryERC20ToDenomResponse> {
    const data = QueryERC20ToDenomRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "ERC20ToDenom", data);
    return promise.then(data => QueryERC20ToDenomResponse.decode(new _m0.Reader(data)));
  }
  denomToERC20(request: QueryDenomToERC20Request): Promise<QueryDenomToERC20Response> {
    const data = QueryDenomToERC20Request.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "DenomToERC20", data);
    return promise.then(data => QueryDenomToERC20Response.decode(new _m0.Reader(data)));
  }
  getLastObservedEthBlock(request: QueryLastObservedEthBlockRequest): Promise<QueryLastObservedEthBlockResponse> {
    const data = QueryLastObservedEthBlockRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetLastObservedEthBlock", data);
    return promise.then(data => QueryLastObservedEthBlockResponse.decode(new _m0.Reader(data)));
  }
  getLastObservedEthNonce(request: QueryLastObservedEthNonceRequest): Promise<QueryLastObservedEthNonceResponse> {
    const data = QueryLastObservedEthNonceRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetLastObservedEthNonce", data);
    return promise.then(data => QueryLastObservedEthNonceResponse.decode(new _m0.Reader(data)));
  }
  getLastObservedERC721EthNonce(request: QueryLastObservedERC721EthNonceRequest = {}): Promise<QueryLastObservedERC721EthNonceResponse> {
    const data = QueryLastObservedERC721EthNonceRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetLastObservedERC721EthNonce", data);
    return promise.then(data => QueryLastObservedERC721EthNonceResponse.decode(new _m0.Reader(data)));
  }
  getAttestations(request: QueryAttestationsRequest): Promise<QueryAttestationsResponse> {
    const data = QueryAttestationsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetAttestations", data);
    return promise.then(data => QueryAttestationsResponse.decode(new _m0.Reader(data)));
  }
  getERC721Attestations(request: QueryERC721AttestationsRequest): Promise<QueryERC721AttestationsResponse> {
    const data = QueryERC721AttestationsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetERC721Attestations", data);
    return promise.then(data => QueryERC721AttestationsResponse.decode(new _m0.Reader(data)));
  }
  getDelegateKeyByValidator(request: QueryDelegateKeysByValidatorAddress): Promise<QueryDelegateKeysByValidatorAddressResponse> {
    const data = QueryDelegateKeysByValidatorAddress.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetDelegateKeyByValidator", data);
    return promise.then(data => QueryDelegateKeysByValidatorAddressResponse.decode(new _m0.Reader(data)));
  }
  getDelegateKeyByEth(request: QueryDelegateKeysByEthAddress): Promise<QueryDelegateKeysByEthAddressResponse> {
    const data = QueryDelegateKeysByEthAddress.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetDelegateKeyByEth", data);
    return promise.then(data => QueryDelegateKeysByEthAddressResponse.decode(new _m0.Reader(data)));
  }
  getDelegateKeyByOrchestrator(request: QueryDelegateKeysByOrchestratorAddress): Promise<QueryDelegateKeysByOrchestratorAddressResponse> {
    const data = QueryDelegateKeysByOrchestratorAddress.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetDelegateKeyByOrchestrator", data);
    return promise.then(data => QueryDelegateKeysByOrchestratorAddressResponse.decode(new _m0.Reader(data)));
  }
  getPendingSendToEth(request: QueryPendingSendToEth): Promise<QueryPendingSendToEthResponse> {
    const data = QueryPendingSendToEth.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetPendingSendToEth", data);
    return promise.then(data => QueryPendingSendToEthResponse.decode(new _m0.Reader(data)));
  }
  getPendingIbcAutoForwards(request: QueryPendingIbcAutoForwards): Promise<QueryPendingIbcAutoForwardsResponse> {
    const data = QueryPendingIbcAutoForwards.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetPendingIbcAutoForwards", data);
    return promise.then(data => QueryPendingIbcAutoForwardsResponse.decode(new _m0.Reader(data)));
  }
  getPendingERC721IbcAutoForwards(request: QueryPendingERC721IbcAutoForwardsRequest): Promise<QueryPendingERC721IbcAutoForwardsResponse> {
    const data = QueryPendingERC721IbcAutoForwardsRequest.encode(request).finish();
    const promise = this.rpc.request("gravity.v1.Query", "GetPendingERC721IbcAutoForwards", data);
    return promise.then(data => QueryPendingERC721IbcAutoForwardsResponse.decode(new _m0.Reader(data)));
  }
}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },
    currentValset(request?: QueryCurrentValsetRequest): Promise<QueryCurrentValsetResponse> {
      return queryService.currentValset(request);
    },
    valsetRequest(request: QueryValsetRequestRequest): Promise<QueryValsetRequestResponse> {
      return queryService.valsetRequest(request);
    },
    valsetConfirm(request: QueryValsetConfirmRequest): Promise<QueryValsetConfirmResponse> {
      return queryService.valsetConfirm(request);
    },
    valsetConfirmsByNonce(request: QueryValsetConfirmsByNonceRequest): Promise<QueryValsetConfirmsByNonceResponse> {
      return queryService.valsetConfirmsByNonce(request);
    },
    lastValsetRequests(request?: QueryLastValsetRequestsRequest): Promise<QueryLastValsetRequestsResponse> {
      return queryService.lastValsetRequests(request);
    },
    lastPendingValsetRequestByAddr(request: QueryLastPendingValsetRequestByAddrRequest): Promise<QueryLastPendingValsetRequestByAddrResponse> {
      return queryService.lastPendingValsetRequestByAddr(request);
    },
    lastPendingBatchRequestByAddr(request: QueryLastPendingBatchRequestByAddrRequest): Promise<QueryLastPendingBatchRequestByAddrResponse> {
      return queryService.lastPendingBatchRequestByAddr(request);
    },
    lastPendingLogicCallByAddr(request: QueryLastPendingLogicCallByAddrRequest): Promise<QueryLastPendingLogicCallByAddrResponse> {
      return queryService.lastPendingLogicCallByAddr(request);
    },
    lastEventNonceByAddr(request: QueryLastEventNonceByAddrRequest): Promise<QueryLastEventNonceByAddrResponse> {
      return queryService.lastEventNonceByAddr(request);
    },
    lastERC721EventNonceByAddr(request: QueryLastERC721EventNonceByAddrRequest): Promise<QueryLastERC721EventNonceByAddrResponse> {
      return queryService.lastERC721EventNonceByAddr(request);
    },
    batchFees(request?: QueryBatchFeeRequest): Promise<QueryBatchFeeResponse> {
      return queryService.batchFees(request);
    },
    outgoingTxBatches(request?: QueryOutgoingTxBatchesRequest): Promise<QueryOutgoingTxBatchesResponse> {
      return queryService.outgoingTxBatches(request);
    },
    outgoingLogicCalls(request?: QueryOutgoingLogicCallsRequest): Promise<QueryOutgoingLogicCallsResponse> {
      return queryService.outgoingLogicCalls(request);
    },
    batchRequestByNonce(request: QueryBatchRequestByNonceRequest): Promise<QueryBatchRequestByNonceResponse> {
      return queryService.batchRequestByNonce(request);
    },
    batchConfirms(request: QueryBatchConfirmsRequest): Promise<QueryBatchConfirmsResponse> {
      return queryService.batchConfirms(request);
    },
    logicConfirms(request: QueryLogicConfirmsRequest): Promise<QueryLogicConfirmsResponse> {
      return queryService.logicConfirms(request);
    },
    eRC20ToDenom(request: QueryERC20ToDenomRequest): Promise<QueryERC20ToDenomResponse> {
      return queryService.eRC20ToDenom(request);
    },
    denomToERC20(request: QueryDenomToERC20Request): Promise<QueryDenomToERC20Response> {
      return queryService.denomToERC20(request);
    },
    getLastObservedEthBlock(request: QueryLastObservedEthBlockRequest): Promise<QueryLastObservedEthBlockResponse> {
      return queryService.getLastObservedEthBlock(request);
    },
    getLastObservedEthNonce(request: QueryLastObservedEthNonceRequest): Promise<QueryLastObservedEthNonceResponse> {
      return queryService.getLastObservedEthNonce(request);
    },
    getLastObservedERC721EthNonce(request?: QueryLastObservedERC721EthNonceRequest): Promise<QueryLastObservedERC721EthNonceResponse> {
      return queryService.getLastObservedERC721EthNonce(request);
    },
    getAttestations(request: QueryAttestationsRequest): Promise<QueryAttestationsResponse> {
      return queryService.getAttestations(request);
    },
    getERC721Attestations(request: QueryERC721AttestationsRequest): Promise<QueryERC721AttestationsResponse> {
      return queryService.getERC721Attestations(request);
    },
    getDelegateKeyByValidator(request: QueryDelegateKeysByValidatorAddress): Promise<QueryDelegateKeysByValidatorAddressResponse> {
      return queryService.getDelegateKeyByValidator(request);
    },
    getDelegateKeyByEth(request: QueryDelegateKeysByEthAddress): Promise<QueryDelegateKeysByEthAddressResponse> {
      return queryService.getDelegateKeyByEth(request);
    },
    getDelegateKeyByOrchestrator(request: QueryDelegateKeysByOrchestratorAddress): Promise<QueryDelegateKeysByOrchestratorAddressResponse> {
      return queryService.getDelegateKeyByOrchestrator(request);
    },
    getPendingSendToEth(request: QueryPendingSendToEth): Promise<QueryPendingSendToEthResponse> {
      return queryService.getPendingSendToEth(request);
    },
    getPendingIbcAutoForwards(request: QueryPendingIbcAutoForwards): Promise<QueryPendingIbcAutoForwardsResponse> {
      return queryService.getPendingIbcAutoForwards(request);
    },
    getPendingERC721IbcAutoForwards(request: QueryPendingERC721IbcAutoForwardsRequest): Promise<QueryPendingERC721IbcAutoForwardsResponse> {
      return queryService.getPendingERC721IbcAutoForwards(request);
    }
  };
};