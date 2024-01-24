/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "andromedad.nibtdfvas";

/** Params defines the parameters for the module. */
export interface Params {
  TokenOutflowPerBlock: number;
  DirectToValidatorPercent: number;
}

const baseParams: object = {
  TokenOutflowPerBlock: 0,
  DirectToValidatorPercent: 0,
};

export const Params = {
  encode(message: Params, writer: Writer = Writer.create()): Writer {
    if (message.TokenOutflowPerBlock !== 0) {
      writer.uint32(8).int64(message.TokenOutflowPerBlock);
    }
    if (message.DirectToValidatorPercent !== 0) {
      writer.uint32(16).int64(message.DirectToValidatorPercent);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseParams } as Params;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.TokenOutflowPerBlock = longToNumber(reader.int64() as Long);
          break;
        case 2:
          message.DirectToValidatorPercent = longToNumber(
            reader.int64() as Long
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Params {
    const message = { ...baseParams } as Params;
    if (
      object.TokenOutflowPerBlock !== undefined &&
      object.TokenOutflowPerBlock !== null
    ) {
      message.TokenOutflowPerBlock = Number(object.TokenOutflowPerBlock);
    } else {
      message.TokenOutflowPerBlock = 0;
    }
    if (
      object.DirectToValidatorPercent !== undefined &&
      object.DirectToValidatorPercent !== null
    ) {
      message.DirectToValidatorPercent = Number(
        object.DirectToValidatorPercent
      );
    } else {
      message.DirectToValidatorPercent = 0;
    }
    return message;
  },

  toJSON(message: Params): unknown {
    const obj: any = {};
    message.TokenOutflowPerBlock !== undefined &&
      (obj.TokenOutflowPerBlock = message.TokenOutflowPerBlock);
    message.DirectToValidatorPercent !== undefined &&
      (obj.DirectToValidatorPercent = message.DirectToValidatorPercent);
    return obj;
  },

  fromPartial(object: DeepPartial<Params>): Params {
    const message = { ...baseParams } as Params;
    if (
      object.TokenOutflowPerBlock !== undefined &&
      object.TokenOutflowPerBlock !== null
    ) {
      message.TokenOutflowPerBlock = object.TokenOutflowPerBlock;
    } else {
      message.TokenOutflowPerBlock = 0;
    }
    if (
      object.DirectToValidatorPercent !== undefined &&
      object.DirectToValidatorPercent !== null
    ) {
      message.DirectToValidatorPercent = object.DirectToValidatorPercent;
    } else {
      message.DirectToValidatorPercent = 0;
    }
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}
