/// Class defines the class of the nft type.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Class {
    /// id defines the unique identifier of the NFT classification, similar to the contract address of ERC721
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
    /// name defines the human-readable name of the NFT classification. Optional
    #[prost(string, tag = "2")]
    pub name: ::prost::alloc::string::String,
    /// symbol is an abbreviated name for nft classification. Optional
    #[prost(string, tag = "3")]
    pub symbol: ::prost::alloc::string::String,
    /// description is a brief description of nft classification. Optional
    #[prost(string, tag = "4")]
    pub description: ::prost::alloc::string::String,
    /// uri for the class metadata stored off chain. It can define schema for Class and NFT `Data` attributes. Optional
    #[prost(string, tag = "5")]
    pub uri: ::prost::alloc::string::String,
    /// uri_hash is a hash of the document pointed by uri. Optional
    #[prost(string, tag = "6")]
    pub uri_hash: ::prost::alloc::string::String,
    /// data is the app specific metadata of the NFT class. Optional
    #[prost(message, optional, tag = "7")]
    pub data: ::core::option::Option<::prost_types::Any>,
}
/// NFT defines the NFT.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Nft {
    /// class_id associated with the NFT, similar to the contract address of ERC721
    #[prost(string, tag = "1")]
    pub class_id: ::prost::alloc::string::String,
    /// id is a unique identifier of the NFT
    #[prost(string, tag = "2")]
    pub id: ::prost::alloc::string::String,
    /// uri for the NFT metadata stored off chain
    #[prost(string, tag = "3")]
    pub uri: ::prost::alloc::string::String,
    /// uri_hash is a hash of the document pointed by uri
    #[prost(string, tag = "4")]
    pub uri_hash: ::prost::alloc::string::String,
    /// data is an app specific data of the NFT. Optional
    #[prost(message, optional, tag = "10")]
    pub data: ::core::option::Option<::prost_types::Any>,
}
/// QueryBalanceRequest is the request type for the Query/Balance RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryBalanceRequest {
    #[prost(string, tag = "1")]
    pub class_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub owner: ::prost::alloc::string::String,
}
/// QueryBalanceResponse is the response type for the Query/Balance RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryBalanceResponse {
    #[prost(uint64, tag = "1")]
    pub amount: u64,
}
/// QueryOwnerRequest is the request type for the Query/Owner RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryOwnerRequest {
    #[prost(string, tag = "1")]
    pub class_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub id: ::prost::alloc::string::String,
}
/// QueryOwnerResponse is the response type for the Query/Owner RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryOwnerResponse {
    #[prost(string, tag = "1")]
    pub owner: ::prost::alloc::string::String,
}
/// QuerySupplyRequest is the request type for the Query/Supply RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QuerySupplyRequest {
    #[prost(string, tag = "1")]
    pub class_id: ::prost::alloc::string::String,
}
/// QuerySupplyResponse is the response type for the Query/Supply RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QuerySupplyResponse {
    #[prost(uint64, tag = "1")]
    pub amount: u64,
}
/// QueryNFTstRequest is the request type for the Query/NFTs RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryNfTsRequest {
    #[prost(string, tag = "1")]
    pub class_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub owner: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "3")]
    pub pagination:
        ::core::option::Option<super::cosmos_sdk_proto::cosmos::base::query::v1beta1::PageRequest>,
}
/// QueryNFTsResponse is the response type for the Query/NFTs RPC methods
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryNfTsResponse {
    #[prost(message, repeated, tag = "1")]
    pub nfts: ::prost::alloc::vec::Vec<Nft>,
    #[prost(message, optional, tag = "2")]
    pub pagination:
        ::core::option::Option<super::cosmos_sdk_proto::cosmos::base::query::v1beta1::PageResponse>,
}
/// QueryNFTRequest is the request type for the Query/NFT RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryNftRequest {
    #[prost(string, tag = "1")]
    pub class_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub id: ::prost::alloc::string::String,
}
/// QueryNFTResponse is the response type for the Query/NFT RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryNftResponse {
    #[prost(message, optional, tag = "1")]
    pub nft: ::core::option::Option<Nft>,
}
/// QueryClassRequest is the request type for the Query/Class RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryClassRequest {
    #[prost(string, tag = "1")]
    pub class_id: ::prost::alloc::string::String,
}
/// QueryClassResponse is the response type for the Query/Class RPC method
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryClassResponse {
    #[prost(message, optional, tag = "1")]
    pub class: ::core::option::Option<Class>,
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
        /// Balance queries the number of NFTs of a given class owned by the owner, same as balanceOf in ERC721
        pub async fn balance(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryBalanceRequest>,
        ) -> Result<tonic::Response<super::QueryBalanceResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cosmos.nft.v1beta1.Query/Balance");
            self.inner.unary(request.into_request(), path, codec).await
        }
        /// Owner queries the owner of the NFT based on its class and id, same as ownerOf in ERC721
        pub async fn owner(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryOwnerRequest>,
        ) -> Result<tonic::Response<super::QueryOwnerResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cosmos.nft.v1beta1.Query/Owner");
            self.inner.unary(request.into_request(), path, codec).await
        }
        /// Supply queries the number of NFTs from the given class, same as totalSupply of ERC721.
        pub async fn supply(
            &mut self,
            request: impl tonic::IntoRequest<super::QuerySupplyRequest>,
        ) -> Result<tonic::Response<super::QuerySupplyResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cosmos.nft.v1beta1.Query/Supply");
            self.inner.unary(request.into_request(), path, codec).await
        }
        /// NFTs queries all NFTs of a given class or owner,choose at least one of the two, similar to tokenByIndex in
        /// ERC721Enumerable
        pub async fn nf_ts(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryNfTsRequest>,
        ) -> Result<tonic::Response<super::QueryNfTsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cosmos.nft.v1beta1.Query/NFTs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        /// NFT queries an NFT based on its class and id.
        pub async fn nft(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryNftRequest>,
        ) -> Result<tonic::Response<super::QueryNftResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cosmos.nft.v1beta1.Query/NFT");
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
}
/// EventSend is emitted on Msg/Send
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EventSend {
    #[prost(string, tag = "1")]
    pub class_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub sender: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub receiver: ::prost::alloc::string::String,
}
/// EventMint is emitted on Mint
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EventMint {
    #[prost(string, tag = "1")]
    pub class_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub owner: ::prost::alloc::string::String,
}
/// EventBurn is emitted on Burn
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EventBurn {
    #[prost(string, tag = "1")]
    pub class_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub owner: ::prost::alloc::string::String,
}
/// GenesisState defines the nft module's genesis state.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    /// class defines the class of the nft type.
    #[prost(message, repeated, tag = "1")]
    pub classes: ::prost::alloc::vec::Vec<Class>,
    #[prost(message, repeated, tag = "2")]
    pub entries: ::prost::alloc::vec::Vec<Entry>,
}
/// Entry Defines all nft owned by a person
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Entry {
    /// owner is the owner address of the following nft
    #[prost(string, tag = "1")]
    pub owner: ::prost::alloc::string::String,
    /// nfts is a group of nfts of the same owner
    #[prost(message, repeated, tag = "2")]
    pub nfts: ::prost::alloc::vec::Vec<Nft>,
}
