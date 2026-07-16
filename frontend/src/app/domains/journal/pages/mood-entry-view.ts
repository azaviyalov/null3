import { Component, computed, inject } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { MoodEntryApi } from "../services/mood-entry-api";
import { MoodEntryDetail } from "../components/mood-entry-detail";
import { map } from "rxjs";
import { toWritableStateSignal } from "../../../core/utils/signal-helpers";
import { stateError, stateSuccess } from "../../../core/utils/state";
import { DiaryEntryLink } from "../models/mood-entry";

@Component({
  selector: "app-mood-entry-view",
  standalone: true,
  imports: [MoodEntryDetail],
  templateUrl: "./mood-entry-view.html",
  styleUrl: "./mood-entry-view.scss",
})
export class MoodEntryView {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly moodEntryApi = inject(MoodEntryApi);

  private readonly entryState = toWritableStateSignal({
    trigger: this.route.params.pipe(map((params) => Number(params["id"]))),
    project: (id) => this.moodEntryApi.getById(id),
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

    this.router.navigate(["/mood/entries", entry.id, "update"]);
  }

  deleteEntry(): void {
    const entry = this.entry();
    if (!entry || entry.deletedAt) {
      this.entryState.set(stateError("Cannot delete a deleted entry"));
      return;
    }

    this.moodEntryApi.delete(entry.id).subscribe({
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

    this.moodEntryApi.restore(entry.id).subscribe({
      next: (restoredEntry) => {
        this.entryState.set(stateSuccess(restoredEntry));
      },
      error: (err) => {
        this.entryState.set(stateError(err));
      },
    });
  }

  openDiaryEntry(entry: DiaryEntryLink): void {
    this.router.navigate(["/diary/entries", entry.id]);
  }

  createDiaryEntry(): void {
    this.router.navigate(["/diary/entries/create"]);
  }

  goBack(): void {
    history.back();
  }
}
