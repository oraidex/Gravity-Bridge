#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EventRegisterInterchainAccount {
    #[prost(string, tag = "1")]
    pub owner: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub connection_id: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub version: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EventSubmitTx {
    #[prost(string, tag = "1")]
    pub owner: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub connection_id: ::prost::alloc::string::String,
    #[prost(message, repeated, tag = "3")]
    pub msgs: ::prost::alloc::vec::Vec<::prost_types::Any>,
}
/// MsgRegisterAccount defines the payload for Msg/RegisterAccount
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgRegisterAccount {
    #[prost(string, tag = "1")]
    pub owner: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub connection_id: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub version: ::prost::alloc::string::String,
}
/// MsgRegisterAccountResponse defines the response for Msg/RegisterAccount
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgRegisterAccountResponse {}
/// MsgSubmitTx defines the payload for Msg/SubmitTx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitTx {
    #[prost(string, tag = "1")]
    pub owner: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub connection_id: ::prost::alloc::string::String,
    /// msgs represents the transactions to be submitted to the host chain
    #[prost(message, repeated, tag = "3")]
    pub msgs: ::prost::alloc::vec::Vec<::prost_types::Any>,
}
/// MsgSubmitTxResponse defines the response for Msg/SubmitTx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitTxResponse {}
/// Generated client implementations.
pub mod msg_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Msg defines the icaauth Msg service.
    #[derive(Debug, Clone)]
    pub struct MsgClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl MsgClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: std::convert::TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> MsgClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_interceptor<F>(inner: T, interceptor: F) -> MsgClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<http::Request<tonic::body::BoxBody>>>::Error:
                Into<StdError> + Send + Sync,
        {
            MsgClient::new(InterceptedService::new(inner, interceptor))
        }

        /// Register defines a rpc handler for MsgRegisterAccount
        pub async fn register_account(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgRegisterAccount>,
        ) -> Result<tonic::Response<super::MsgRegisterAccountResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gaia.icaauth.v1.Msg/RegisterAccount");
            self.inner.unary(request.into_request(), path, codec).await
        }
        /// SubmitTx defines a rpc handler for MsgSubmitTx
        pub async fn submit_tx(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgSubmitTx>,
        ) -> Result<tonic::Response<super::MsgSubmitTxResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gaia.icaauth.v1.Msg/SubmitTx");
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
}
/// QueryInterchainAccountFromAddressRequest is the request type for the Query/InterchainAccountAddress RPC
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryInterchainAccountFromAddressRequest {
    #[prost(string, tag = "1")]
    pub owner: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub connection_id: ::prost::alloc::string::String,
}
/// QueryInterchainAccountFromAddressResponse the response type for the Query/InterchainAccountFromAddress RPC
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryInterchainAccountFromAddressResponse {
    #[prost(string, tag = "1")]
    pub interchain_account_address: ::prost::alloc::string::String,
}
/// QueryInterchainAccountsWithConnectionRequest is the request type for the Query/InterchainAccountsWithConnection RPC
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryInterchainAccountsWithConnectionRequest {
    #[prost(string, tag = "1")]
    pub connection_id: ::prost::alloc::string::String,
}
/// QueryInterchainAccountFromAddressResponse the response type for the Query/InterchainAccountsWithConnection RPC
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryInterchainAccountsWithConnectionResponse {
    #[prost(message, repeated, tag = "1")]
    pub interchain_accounts: ::prost::alloc::vec::Vec<
        cosmos_sdk_proto::ibc::applications::interchain_accounts::v1::RegisteredInterchainAccount,
    >,
}
/// Generated client implementations.
pub mod query_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Query defines the gRPC querier service.
    #[derive(Debug, Clone)]
    pub struct QueryClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl QueryClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: std::convert::TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> QueryClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> QueryClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<http::Request<tonic::body::BoxBody>>>::Error:
                Into<StdError> + Send + Sync,
        {
            QueryClient::new(InterceptedService::new(inner, interceptor))
        }

        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }

        /// QueryInterchainAccountFromAddress returns the interchain account for given owner address on a given connection
        pub async fn interchain_account_from_address(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryInterchainAccountFromAddressRequest>,
        ) -> Result<tonic::Response<super::QueryInterchainAccountFromAddressResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/gaia.icaauth.v1.Query/InterchainAccountFromAddress",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        /// QueryInterchainAccountsWithConnection returns all the interchain accounts on a given connection
        pub async fn interchain_accounts_with_connection(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryInterchainAccountsWithConnectionRequest>,
        ) -> Result<
            tonic::Response<super::QueryInterchainAccountsWithConnectionResponse>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/gaia.icaauth.v1.Query/InterchainAccountsWithConnection",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
}
