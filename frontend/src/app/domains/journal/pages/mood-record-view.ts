import { Component, DestroyRef, computed, inject, signal } from "@angular/core";
import {
  rxResource,
  takeUntilDestroyed,
  toSignal,
} from "@angular/core/rxjs-interop";
import { ActivatedRoute, Router } from "@angular/router";
import { map, type Observable } from "rxjs";
import { MoodRecordApi } from "../services/mood-record-api";
import { MoodRecordDetail } from "../components/mood-record-detail";
import type { DiaryEntryLink, MoodRecord } from "../models/mood-record";

@Component({
  selector: "app-mood-record-view",
  standalone: true,
  imports: [MoodRecordDetail],
  templateUrl: "./mood-record-view.html",
  styleUrl: "./mood-record-view.scss",
})
export class MoodRecordView {
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

    void this.#router.navigate(["/mood-records", entry.id, "update"]);
  }

  protected deleteEntry(): void {
    const entry = this.entry();
    if (!entry || entry.deletedAt) {
      this.#actionError.set("Cannot delete a deleted entry");
      return;
    }

    this.#applyMutation(this.#moodRecordApi.delete(entry.id));
  }

  protected restoreEntry(): void {
    const entry = this.entry();
    if (!entry?.deletedAt) {
      this.#actionError.set("Cannot restore an active entry");
      return;
    }

    this.#applyMutation(this.#moodRecordApi.restore(entry.id));
  }

  protected openDiaryEntry(entry: DiaryEntryLink): void {
    void this.#router.navigate(["/diary-entries", entry.id]);
  }

  protected createDiaryEntry(): void {
    void this.#router.navigate(["/diary-entries/create"]);
  }

  protected goBack(): void {
    history.back();
  }

  #applyMutation(request: Observable<MoodRecord>): void {
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
