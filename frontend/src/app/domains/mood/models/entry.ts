export class Entry {
  constructor(
    readonly id: number,
    readonly feeling: string,
    readonly createdAt: Date,
    readonly updatedAt: Date,
    readonly deletedAt?: Date,
    readonly note?: string,
  ) {}

  static fromResponse(data: EntryResponse): Entry {
    return new Entry(
      data.id,
      data.feeling,
      new Date(data.created_at),
      new Date(data.updated_at),
      data.deleted_at ? new Date(data.deleted_at) : undefined,
      data.note || undefined,
    );
  }
}

export interface EditEntryRequest {
  readonly feeling: string;
  readonly note?: string;
}

export interface EntryResponse {
  readonly id: number;
  readonly feeling: string;
  readonly user_id: number;
  readonly note?: string;
  readonly created_at: string;
  readonly updated_at: string;
  readonly deleted_at?: string;
}
