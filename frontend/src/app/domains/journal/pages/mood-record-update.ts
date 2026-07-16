import { Component, computed, inject, signal } from "@angular/core";
import { EditMoodRecordRequest, MoodRecord } from "../models/mood-record";
import { ActivatedRoute, Router } from "@angular/router";
import { MoodRecordApi } from "../services/mood-record-api";
import { MoodRecordForm } from "../components/mood-record-form";
import { map } from "rxjs";
import { toWritableStateSignal } from "../../../core/utils/signal-helpers";
import { stateError } from "../../../core/utils/state";

@Component({
  selector: "app-mood-record-update",
  standalone: true,
  imports: [MoodRecordForm],
  templateUrl: "./mood-record-update.html",
  styleUrl: "./mood-record-update.scss",
})
export class MoodRecordUpdate {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly moodRecordApi = inject(MoodRecordApi);

  readonly isSubmitting = signal(false);

  private readonly entryState = toWritableStateSignal({
    trigger: this.route.params.pipe(map((params) => Number(params["id"]))),
    project: (id) => this.moodRecordApi.getById(id),
  });

  readonly entry = computed(() => this.entryState().value);
  readonly isLoading = computed(() => this.entryState().isLoading);
  readonly errorMessage = computed(() => this.entryState().error);

  submit(payload: EditMoodRecordRequest): void {
    const entry = this.entry();
    if (!entry) {
      return;
    }

    this.isSubmitting.set(true);

    this.moodRecordApi.update(entry.id, payload).subscribe({
      next: (entry: MoodRecord) => this.handleSuccess(entry),
      error: (err) => this.handleError(err),
    });
  }

  private handleSuccess(entry: MoodRecord): void {
    this.isSubmitting.set(false);
    this.router.navigate(["/mood/records", entry.id]);
  }

  private handleError(err: unknown): void {
    this.isSubmitting.set(false);
    this.entryState.set(stateError(err));
    console.error(err);
  }
}
