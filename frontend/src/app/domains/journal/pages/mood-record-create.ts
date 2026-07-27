import { Component, DestroyRef, inject, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import type { EditMoodRecordRequest } from "../models/mood-record";
import { MoodRecordApi } from "../services/mood-record-api";
import { Router } from "@angular/router";
import { finalize } from "rxjs";
import { MoodRecordForm } from "../components/mood-record-form";

@Component({
  selector: "app-mood-record-create",
  standalone: true,
  imports: [MoodRecordForm],
  templateUrl: "./mood-record-create.html",
  styleUrl: "./mood-record-create.scss",
})
export class MoodRecordCreate {
  readonly #router = inject(Router);
  readonly #moodRecordApi = inject(MoodRecordApi);
  readonly #destroyRef = inject(DestroyRef);

  protected readonly isSubmitting = signal(false);
  protected readonly errorMessage = signal<string | null>(null);

  protected submit(payload: EditMoodRecordRequest): void {
    this.isSubmitting.set(true);
    this.errorMessage.set(null);

    this.#moodRecordApi
      .create(payload)
      .pipe(
        takeUntilDestroyed(this.#destroyRef),
        finalize(() => this.isSubmitting.set(false)),
      )
      .subscribe({
        next: (entry) => {
          void this.#router.navigate(["/mood-records", entry.id]);
        },
        error: (error: unknown) => {
          this.errorMessage.set("Could not create the mood record.");
          console.error(error);
        },
      });
  }
}
