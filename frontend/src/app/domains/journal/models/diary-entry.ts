import { MoodEntry, MoodEntryResponse } from "./mood-entry";

export class DiaryEntry {
  constructor(
    readonly id: number,
    readonly title: string | undefined,
    readonly markdown: string,
    readonly preview: string | undefined,
    readonly occurredAt: Date,
    readonly createdAt: Date,
    readonly updatedAt: Date,
    readonly deletedAt: Date | undefined,
    readonly referencedMoodEntries: MoodEntry[],
  ) {}

  get headline(): string {
    return this.title || this.preview || "Untitled diary entry";
  }

  get excerpt(): string {
    return this.preview || "No text preview.";
  }

  static fromResponse(data: DiaryEntryResponse): DiaryEntry {
    return new DiaryEntry(
      data.id,
      data.title || undefined,
      data.markdown,
      data.preview || undefined,
      new Date(data.occurred_at),
      new Date(data.created_at),
      new Date(data.updated_at),
      data.deleted_at ? new Date(data.deleted_at) : undefined,
      (data.referenced_mood_entries ?? []).map(MoodEntry.fromResponse),
    );
  }
}

export interface EditDiaryEntryRequest {
  readonly title?: string;
  readonly markdown: string;
  readonly occurred_at: string;
}

export interface DiaryEntryResponse {
  readonly id: number;
  readonly user_id: number;
  readonly title?: string;
  readonly markdown: string;
  readonly preview?: string;
  readonly occurred_at: string;
  readonly created_at: string;
  readonly updated_at: string;
  readonly deleted_at: string | null;
  readonly referenced_mood_entries?: MoodEntryResponse[];
}
