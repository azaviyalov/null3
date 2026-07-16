import { Component, computed, inject, signal } from "@angular/core";
import { EditMoodEntryRequest, MoodEntry } from "../models/mood-entry";
import { ActivatedRoute, Router } from "@angular/router";
import { MoodEntryApi } from "../services/mood-entry-api";
import { MoodEntryForm } from "../components/mood-entry-form";
import { map } from "rxjs";
import { toWritableStateSignal } from "../../../core/utils/signal-helpers";
import { stateError } from "../../../core/utils/state";

@Component({
  selector: "app-mood-entry-update",
  standalone: true,
  imports: [MoodEntryForm],
  templateUrl: "./mood-entry-update.html",
  styleUrl: "./mood-entry-update.scss",
})
export class MoodEntryUpdate {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly moodEntryApi = inject(MoodEntryApi);

  readonly isSubmitting = signal(false);

  private readonly entryState = toWritableStateSignal({
    trigger: this.route.params.pipe(map((params) => Number(params["id"]))),
    project: (id) => this.moodEntryApi.getById(id),
  });

  readonly entry = computed(() => this.entryState().value);
  readonly isLoading = computed(() => this.entryState().isLoading);
  readonly errorMessage = computed(() => this.entryState().error);

  submit(payload: EditMoodEntryRequest): void {
    const entry = this.entry();
    if (!entry) {
      return;
    }

    this.isSubmitting.set(true);

    this.moodEntryApi.update(entry.id, payload).subscribe({
      next: (entry: MoodEntry) => this.handleSuccess(entry),
      error: (err) => this.handleError(err),
    });
  }

  private handleSuccess(entry: MoodEntry): void {
    this.isSubmitting.set(false);
    this.router.navigate(["/mood/entries", entry.id]);
  }

  private handleError(err: unknown): void {
    this.isSubmitting.set(false);
    this.entryState.set(stateError(err));
    console.error(err);
  }
}
