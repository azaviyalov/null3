import { Component, DestroyRef, computed, inject, signal } from "@angular/core";
import {
  rxResource,
  takeUntilDestroyed,
  toSignal,
} from "@angular/core/rxjs-interop";
import type { EditMoodRecordRequest } from "../models/mood-record";
import { ActivatedRoute, Router } from "@angular/router";
import { finalize, map } from "rxjs";
import { MoodRecordApi } from "../services/mood-record-api";
import { MoodRecordForm } from "../components/mood-record-form";

@Component({
  selector: "app-mood-record-update",
  standalone: true,
  imports: [MoodRecordForm],
  templateUrl: "./mood-record-update.html",
  styleUrl: "./mood-record-update.scss",
})
export class MoodRecordUpdate {
  readonly #route = inject(ActivatedRoute);
  readonly #router = inject(Router);
  readonly #moodRecordApi = inject(MoodRecordApi);
  readonly #destroyRef = inject(DestroyRef);

  readonly #entryId = toSignal(
    this.#route.paramMap.pipe(map((params) => Number(params.get("id")))),
  );
  readonly #entryResource = rxResource({
    params: () => this.#entryId(),
    stream: ({ params: id }) => this.#moodRecordApi.getById(id),
  });
  readonly #submissionError = signal<string | null>(null);

  protected readonly entry = computed(() =>
    this.#entryResource.hasValue() ? this.#entryResource.value() : null,
  );
  protected readonly isLoading = this.#entryResource.isLoading;
  protected readonly errorMessage = computed(
    () =>
      this.#submissionError() ?? this.#entryResource.error()?.message ?? null,
  );
  protected readonly isSubmitting = signal(false);

  protected submit(payload: EditMoodRecordRequest): void {
    const entry = this.entry();
    if (!entry) {
      return;
    }

    this.isSubmitting.set(true);
    this.#submissionError.set(null);

    this.#moodRecordApi
      .update(entry.id, payload)
      .pipe(
        takeUntilDestroyed(this.#destroyRef),
        finalize(() => this.isSubmitting.set(false)),
      )
      .subscribe({
        next: (updatedEntry) => {
          void this.#router.navigate(["/mood-records", updatedEntry.id]);
        },
        error: (error: unknown) => {
          this.#submissionError.set(
            error instanceof Error
              ? error.message
              : "An unknown error occurred",
          );
          console.error(error);
        },
      });
  }
}
