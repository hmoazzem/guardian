import * as jspb from 'google-protobuf'



export class CPUClock extends jspb.Message {
  getCpuClockList(): Array<number>;
  setCpuClockList(value: Array<number>): CPUClock;
  clearCpuClockList(): CPUClock;
  addCpuClock(value: number, index?: number): CPUClock;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CPUClock.AsObject;
  static toObject(includeInstance: boolean, msg: CPUClock): CPUClock.AsObject;
  static serializeBinaryToWriter(message: CPUClock, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CPUClock;
  static deserializeBinaryFromReader(message: CPUClock, reader: jspb.BinaryReader): CPUClock;
}

export namespace CPUClock {
  export type AsObject = {
    cpuClockList: Array<number>,
  }
}

export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  }
}

