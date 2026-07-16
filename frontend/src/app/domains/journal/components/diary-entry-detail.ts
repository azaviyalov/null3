import { CommonModule } from "@angular/common";
import { Component, computed, input, output } from "@angular/core";
import { DiaryEntry } from "../models/diary-entry";
import { MoodRecord } from "../models/mood-record";
import { MoodRecordHistory } from "./mood-record-history";
import { MarkdownRenderer } from "./markdown-renderer";

@Component({
  selector: "app-diary-entry-detail",
  standalone: true,
  imports: [CommonModule, MarkdownRenderer, MoodRecordHistory],
  templateUrl: "./diary-entry-detail.html",
  styleUrl: "./diary-entry-detail.scss",
})
export class DiaryEntryDetail {
  readonly skeleton = input(false);
  readonly entry = input<DiaryEntry | null>(null);
  readonly showEdit = input(false);
  readonly showDelete = input(false);
  readonly showRestore = input(false);

  readonly showActions = computed(
    () => this.showEdit() || this.showDelete() || this.showRestore(),
  );
  readonly showFooter = computed(
    () =>
      this.entry() &&
      (this.entry()?.deletedAt ||
        this.entry()?.createdAt?.getTime() !==
          this.entry()?.updatedAt?.getTime()) &&
      !this.skeleton(),
  );
  readonly updated = computed(
    () =>
      this.entry() &&
      this.entry()?.updatedAt.getTime() !== this.entry()?.createdAt.getTime(),
  );
  readonly deleted = computed(() => !!this.entry()?.deletedAt);

  readonly edit = output<void>();
  readonly delete = output<void>();
  readonly restore = output<void>();
  readonly openMood = output<MoodRecord>();
}
