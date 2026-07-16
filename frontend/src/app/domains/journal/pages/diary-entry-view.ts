import { Component, computed, inject } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { map } from "rxjs";
import { MoodRecord } from "../models/mood-record";
import { toWritableStateSignal } from "../../../core/utils/signal-helpers";
import { stateError, stateSuccess } from "../../../core/utils/state";
import { DiaryEntryApi } from "../services/diary-entry-api";
import { DiaryEntryDetail } from "../components/diary-entry-detail";

@Component({
  selector: "app-diary-entry-view",
  standalone: true,
  imports: [DiaryEntryDetail],
  templateUrl: "./diary-entry-view.html",
  styleUrl: "./diary-entry-view.scss",
})
export class DiaryEntryView {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly diaryEntryApi = inject(DiaryEntryApi);

  private readonly entryState = toWritableStateSignal({
    trigger: this.route.params.pipe(map((params) => Number(params["id"]))),
    project: (id) => this.diaryEntryApi.getById(id),
  });

  readonly entry = computed(() => this.entryState().value);
  readonly isLoading = computed(() => this.entryState().isLoading);
  readonly errorMessage = computed(() => this.entryState().error);

  editEntry(): void {
    const entry = this.entry();
    if (!entry || entry.deletedAt) {
      this.entryState.set(stateError("Cannot edit a deleted entry"));
      return;
    }

    this.router.navigate(["/diary-entries", entry.id, "update"]);
  }

  deleteEntry(): void {
    const entry = this.entry();
    if (!entry || entry.deletedAt) {
      this.entryState.set(stateError("Cannot delete a deleted entry"));
      return;
    }

    this.diaryEntryApi.delete(entry.id).subscribe({
      next: (deletedEntry) => this.entryState.set(stateSuccess(deletedEntry)),
      error: (err) => this.entryState.set(stateError(err)),
    });
  }

  restoreEntry(): void {
    const entry = this.entry();
    if (!entry?.deletedAt) {
      this.entryState.set(stateError("Cannot restore an active entry"));
      return;
    }

    this.diaryEntryApi.restore(entry.id).subscribe({
      next: (restoredEntry) => this.entryState.set(stateSuccess(restoredEntry)),
      error: (err) => this.entryState.set(stateError(err)),
    });
  }

  openMoodRecord(entry: MoodRecord): void {
    this.router.navigate(["/mood-records", entry.id]);
  }

  goBack(): void {
    history.back();
  }
}
