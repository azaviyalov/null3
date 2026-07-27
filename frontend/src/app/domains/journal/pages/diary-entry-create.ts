import { Component, DestroyRef, inject, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { Router } from "@angular/router";
import { finalize } from "rxjs";
import { DiaryEntryApi } from "../services/diary-entry-api";
import type { EditDiaryEntryRequest } from "../models/diary-entry";
import { DiaryEntryForm } from "../components/diary-entry-form";

@Component({
  selector: "app-diary-entry-create",
  standalone: true,
  imports: [DiaryEntryForm],
  templateUrl: "./diary-entry-create.html",
  styleUrl: "./diary-entry-create.scss",
})
export class DiaryEntryCreate {
  readonly #router = inject(Router);
  readonly #diaryEntryApi = inject(DiaryEntryApi);
  readonly #destroyRef = inject(DestroyRef);

  protected readonly isSubmitting = signal(false);
  protected readonly errorMessage = signal<string | null>(null);

  protected submit(payload: EditDiaryEntryRequest): void {
    this.isSubmitting.set(true);
    this.errorMessage.set(null);

    this.#diaryEntryApi
      .create(payload)
      .pipe(
        takeUntilDestroyed(this.#destroyRef),
        finalize(() => this.isSubmitting.set(false)),
      )
      .subscribe({
        next: (entry) => {
          void this.#router.navigate(["/diary-entries", entry.id]);
        },
        error: (error: unknown) => {
          this.errorMessage.set("Could not create the diary entry.");
          console.error(error);
        },
      });
  }
}
