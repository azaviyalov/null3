export class DiaryEntryLink {
  constructor(
    readonly id: number,
    readonly title: string | undefined,
    readonly preview: string | undefined,
    readonly occurredAt: Date,
    readonly createdAt: Date,
    readonly updatedAt: Date,
  ) {}

  get headline(): string {
    return this.title || this.preview || "Untitled diary entry";
  }

  static fromResponse(data: DiaryEntryLinkResponse): DiaryEntryLink {
    return new DiaryEntryLink(
      data.id,
      data.title || undefined,
      data.preview || undefined,
      new Date(data.occurred_at),
      new Date(data.created_at),
      new Date(data.updated_at),
    );
  }
}

export class MoodEntry {
  constructor(
    readonly id: number,
    readonly feeling: string,
    readonly emoji: string | undefined,
    readonly createdAt: Date,
    readonly updatedAt: Date,
    readonly deletedAt?: Date,
    readonly note?: string,
    readonly diaryEntryLinks: DiaryEntryLink[] = [],
  ) {}

  static fromResponse(data: MoodEntryResponse): MoodEntry {
    return new MoodEntry(
      data.id,
      data.feeling,
      data.emoji || undefined,
      new Date(data.created_at),
      new Date(data.updated_at),
      data.deleted_at ? new Date(data.deleted_at) : undefined,
      data.note || undefined,
      (data.diary_entry_links ?? []).map(DiaryEntryLink.fromResponse),
    );
  }
}

export interface EditMoodEntryRequest {
  readonly feeling: string;
  readonly emoji?: string;
  readonly note?: string;
}

export interface MoodEntryResponse {
  readonly id: number;
  readonly feeling: string;
  readonly emoji?: string;
  readonly user_id: number;
  readonly note?: string;
  readonly created_at: string;
  readonly updated_at: string;
  readonly deleted_at: string | null;
  readonly diary_entry_links?: DiaryEntryLinkResponse[];
}

export interface DiaryEntryLinkResponse {
  readonly id: number;
  readonly title?: string;
  readonly preview?: string;
  readonly occurred_at: string;
  readonly created_at: string;
  readonly updated_at: string;
}
