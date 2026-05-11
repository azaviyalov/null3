import { CommonModule } from "@angular/common";
import { Component, computed, input, output } from "@angular/core";
import { Entry as MoodEntry } from "../../../mood/models/entry";
import { EntryHistory as MoodEntryHistory } from "../../../mood/components/entry-history/entry-history";
import { DiaryEntry } from "../../models/entry";
import { MarkdownRenderer } from "../markdown-renderer/markdown-renderer";

@Component({
  selector: "app-diary-entry-detail",
  standalone: true,
  imports: [CommonModule, MarkdownRenderer, MoodEntryHistory],
  templateUrl: "./entry-detail.html",
  styleUrl: "./entry-detail.scss",
})
export class EntryDetail {
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
  readonly openMood = output<MoodEntry>();
}
