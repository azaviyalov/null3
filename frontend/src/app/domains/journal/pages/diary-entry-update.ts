import { Component, DestroyRef, computed, inject, signal } from "@angular/core";
import {
  rxResource,
  takeUntilDestroyed,
  toSignal,
} from "@angular/core/rxjs-interop";
import { ActivatedRoute, Router } from "@angular/router";
import { finalize, map } from "rxjs";
import type { EditDiaryEntryRequest } from "../models/diary-entry";
import { DiaryEntryApi } from "../services/diary-entry-api";
import { DiaryEntryForm } from "../components/diary-entry-form";

@Component({
  selector: "app-diary-entry-update",
  standalone: true,
  imports: [DiaryEntryForm],
  templateUrl: "./diary-entry-update.html",
  styleUrl: "./diary-entry-update.scss",
})
export class DiaryEntryUpdate {
  readonly #route = inject(ActivatedRoute);
  readonly #router = inject(Router);
  readonly #diaryEntryApi = inject(DiaryEntryApi);
  readonly #destroyRef = inject(DestroyRef);

  readonly #entryId = toSignal(
    this.#route.paramMap.pipe(map((params) => Number(params.get("id")))),
  );
  readonly #entryResource = rxResource({
    params: () => this.#entryId(),
    stream: ({ params: id }) => this.#diaryEntryApi.getById(id),
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

  protected submit(payload: EditDiaryEntryRequest): void {
    const entry = this.entry();
    if (!entry) {
      return;
    }

    this.isSubmitting.set(true);
    this.#submissionError.set(null);

    this.#diaryEntryApi
      .update(entry.id, payload)
      .pipe(
        takeUntilDestroyed(this.#destroyRef),
        finalize(() => this.isSubmitting.set(false)),
      )
      .subscribe({
        next: (updatedEntry) => {
          void this.#router.navigate(["/diary-entries", updatedEntry.id]);
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
