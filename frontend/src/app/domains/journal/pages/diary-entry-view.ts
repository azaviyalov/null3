import { Component, DestroyRef, computed, inject, signal } from "@angular/core";
import {
  rxResource,
  takeUntilDestroyed,
  toSignal,
} from "@angular/core/rxjs-interop";
import { ActivatedRoute, Router } from "@angular/router";
import { map, type Observable } from "rxjs";
import type { DiaryEntry } from "../models/diary-entry";
import type { MoodRecord } from "../models/mood-record";
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
  readonly #actionError = signal<string | null>(null);

  protected readonly entry = computed(() =>
    this.#entryResource.hasValue() ? this.#entryResource.value() : null,
  );
  protected readonly isLoading = this.#entryResource.isLoading;
  protected readonly errorMessage = computed(
    () => this.#actionError() ?? this.#entryResource.error()?.message ?? null,
  );

  protected editEntry(): void {
    const entry = this.entry();
    if (!entry || entry.deletedAt) {
      this.#actionError.set("Cannot edit a deleted entry");
      return;
    }

    void this.#router.navigate(["/diary-entries", entry.id, "update"]);
  }

  protected deleteEntry(): void {
    const entry = this.entry();
    if (!entry || entry.deletedAt) {
      this.#actionError.set("Cannot delete a deleted entry");
      return;
    }

    this.#applyMutation(this.#diaryEntryApi.delete(entry.id));
  }

  protected restoreEntry(): void {
    const entry = this.entry();
    if (!entry?.deletedAt) {
      this.#actionError.set("Cannot restore an active entry");
      return;
    }

    this.#applyMutation(this.#diaryEntryApi.restore(entry.id));
  }

  protected openMoodRecord(entry: MoodRecord): void {
    void this.#router.navigate(["/mood-records", entry.id]);
  }

  protected goBack(): void {
    history.back();
  }

  #applyMutation(request: Observable<DiaryEntry>): void {
    this.#actionError.set(null);

    request.pipe(takeUntilDestroyed(this.#destroyRef)).subscribe({
      next: (entry) => this.#entryResource.set(entry),
      error: (error: unknown) => {
        this.#actionError.set(
          error instanceof Error ? error.message : "An unknown error occurred",
        );
      },
    });
  }
}
