import { Component, computed, input, output } from "@angular/core";
import { CommonModule } from "@angular/common";
import type { MoodRecord } from "../models/mood-record";
import { type DateGroup, groupEntriesByDate } from "../utils/date-grouping";
import { feelingLabel } from "../utils/feeling-presenter";

@Component({
  selector: "app-mood-record-history",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./mood-record-history.html",
  styleUrl: "./mood-record-history.scss",
})
export class MoodRecordHistory {
  readonly skeleton = input(false);
  readonly skeletonCount = input(10);
  readonly entries = input<MoodRecord[] | null>(null);
  readonly showOpen = input(false);
  readonly showEdit = input(false);
  readonly showDelete = input(false);
  readonly showDeletedStubs = input(false);
  readonly emptyMessage = input("No mood records found.");

  protected readonly groupedEntries = computed<DateGroup<MoodRecord>[]>(() =>
    groupEntriesByDate(
      (this.entries() ?? []).filter(
        (entry) => !entry.deletedAt || this.showDeletedStubs(),
      ),
      (entry) => entry.createdAt,
    ),
  );
  protected readonly skeletonRows = computed(() =>
    Array.from({ length: this.skeletonCount() }, (_, index) => index),
  );

  readonly open = output<MoodRecord>();
  readonly edit = output<MoodRecord>();
  readonly delete = output<MoodRecord>();

  protected readonly showActions = computed(
    () => this.showOpen() || this.showEdit() || this.showDelete(),
  );
  protected readonly feelingLabel = feelingLabel;
}
