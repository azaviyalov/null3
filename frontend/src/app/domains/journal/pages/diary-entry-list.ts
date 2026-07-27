import { Component, computed, inject, signal } from "@angular/core";
import { rxResource } from "@angular/core/rxjs-interop";
import { Router } from "@angular/router";
import { DiaryEntryApi } from "../services/diary-entry-api";
import type { DiaryEntry } from "../models/diary-entry";
import { DiaryEntryFeed } from "../components/diary-entry-feed";

@Component({
  selector: "app-diary-entry-list",
  standalone: true,
  imports: [DiaryEntryFeed],
  templateUrl: "./diary-entry-list.html",
  styleUrl: "./diary-entry-list.scss",
})
export class DiaryEntryList {
  readonly #router = inject(Router);
  readonly #diaryEntryApi = inject(DiaryEntryApi);

  protected readonly defaultPageSize = 10;

  readonly #pageOffset = signal(0);
  readonly #pageResource = rxResource({
    params: () => ({
      limit: this.defaultPageSize,
      offset: this.#pageOffset(),
    }),
    stream: ({ params }) =>
      this.#diaryEntryApi.getPaged(params.limit, params.offset),
  });
  protected readonly page = computed(() =>
    this.#pageResource.hasValue() ? this.#pageResource.value() : null,
  );
  protected readonly isLoading = this.#pageResource.isLoading;
  protected readonly canGoPrevious = computed(() => this.#pageOffset() > 0);
  protected readonly canGoNext = computed(
    () =>
      this.#pageResource.hasValue() &&
      this.#pageOffset() + this.defaultPageSize <
        this.#pageResource.value().totalCount,
  );

  protected openEntry(entry: DiaryEntry): void {
    void this.#router.navigate(["/diary-entries", entry.id]);
  }

  protected createEntry(): void {
    void this.#router.navigate(["/diary-entries/create"]);
  }

  protected goToPreviousPage(): void {
    if (this.canGoPrevious()) {
      this.#pageOffset.update((offset) =>
        Math.max(0, offset - this.defaultPageSize),
      );
    }
  }

  protected goToNextPage(): void {
    if (this.canGoNext()) {
      this.#pageOffset.update((offset) => offset + this.defaultPageSize);
    }
  }
}
