export class Entry {
  constructor(
    readonly id: number,
    readonly feeling: string,
    readonly createdAt: Date,
    readonly updatedAt: Date,
    readonly deletedAt?: Date,
    readonly note?: string,
  ) {}

  static fromView(data: EntryView): Entry {
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

export class PaginatedEntryList {
  constructor(
    readonly items: Entry[],
    readonly totalCount: number,
  ) {}

  static fromView(data: PaginatedEntryListView): PaginatedEntryList {
    return new PaginatedEntryList(
      data.items.map(Entry.fromView),
      data.total_count,
    );
  }

  static empty(): PaginatedEntryList {
    return new PaginatedEntryList([], 0);
  }
}

export interface EditEntryRequest {
  readonly feeling: string;
  readonly note?: string;
}

export interface EntryView {
  readonly id: number;
  readonly feeling: string;
  readonly user_id: number;
  readonly note?: string;
  readonly created_at: string;
  readonly updated_at: string;
  readonly deleted_at?: string;
}

export interface PaginatedEntryListView {
  readonly items: EntryView[];
  readonly total_count: number;
}
