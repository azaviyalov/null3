import { CommonModule } from "@angular/common";
import { Component, computed, input, output } from "@angular/core";
import type { DiaryEntry } from "../models/diary-entry";
import { type DateGroup, groupEntriesByDate } from "../utils/date-grouping";

@Component({
  selector: "app-diary-entry-feed",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./diary-entry-feed.html",
  styleUrl: "./diary-entry-feed.scss",
})
export class DiaryEntryFeed {
  readonly skeleton = input(false);
  readonly skeletonCount = input(10);
  readonly entries = input<DiaryEntry[] | null>(null);
  readonly showOpen = input(false);
  readonly showEdit = input(false);
  readonly showDelete = input(false);
  readonly emptyMessage = input("No diary entries found.");

  protected readonly groupedEntries = computed<DateGroup<DiaryEntry>[]>(() =>
    groupEntriesByDate(
      (this.entries() ?? []).filter((entry) => !entry.deletedAt),
      (entry) => entry.occurredAt,
    ),
  );
  protected readonly skeletonRows = computed(() =>
    Array.from({ length: this.skeletonCount() }, (_, index) => index),
  );
  protected readonly showActions = computed(
    () => this.showOpen() || this.showEdit() || this.showDelete(),
  );

  readonly open = output<DiaryEntry>();
  readonly edit = output<DiaryEntry>();
  readonly delete = output<DiaryEntry>();
}
