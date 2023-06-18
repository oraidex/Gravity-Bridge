import * as _135 from "./v1/attestation";
import * as _136 from "./v1/batch";
import * as _137 from "./v1/ethereum_signer";
import * as _138 from "./v1/genesis";
import * as _139 from "./v1/msgs";
import * as _140 from "./v1/pool";
import * as _141 from "./v1/query";
import * as _142 from "./v1/types";
import * as _227 from "./v1/msgs.amino";
import * as _228 from "./v1/msgs.registry";
import * as _229 from "./v1/query.rpc.Query";
import * as _230 from "./v1/msgs.rpc.msg";
import * as _233 from "./rpc.query";
import * as _234 from "./rpc.tx";
export namespace gravity {
  export const v1 = {
    ..._135,
    ..._136,
    ..._137,
    ..._138,
    ..._139,
    ..._140,
    ..._141,
    ..._142,
    ..._227,
    ..._228,
    ..._229,
    ..._230
  };
  export const ClientFactory = {
    ..._233,
    ..._234
  };
}