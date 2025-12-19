export class Page<T> {
  constructor(
    readonly items: T[],
    readonly totalCount: number,
  ) {}

  static fromResponse<TOut, TRaw>(
    data: PageResponse<TRaw>,
    mapItem: (raw: TRaw) => TOut,
  ): Page<TOut> {
    return new Page<TOut>(data.items.map(mapItem), data.total_count);
  }
}

export interface PageResponse<T> {
  readonly items: T[];
  readonly total_count: number;
}
