import { Component, computed, inject, signal } from "@angular/core";
import { rxResource } from "@angular/core/rxjs-interop";
import type { MoodRecord } from "../models/mood-record";
import { MoodRecordApi } from "../services/mood-record-api";
import { Router } from "@angular/router";
import { MoodRecordHistory } from "../components/mood-record-history";

@Component({
  selector: "app-mood-record-list",
  standalone: true,
  imports: [MoodRecordHistory],
  templateUrl: "./mood-record-list.html",
  styleUrl: "./mood-record-list.scss",
})
export class MoodRecordList {
  readonly #router = inject(Router);
  readonly #moodRecordApi = inject(MoodRecordApi);

  protected readonly defaultPageSize = 10;

  readonly #pageOffset = signal(0);
  readonly #pageResource = rxResource({
    params: () => ({
      limit: this.defaultPageSize,
      offset: this.#pageOffset(),
    }),
    stream: ({ params }) =>
      this.#moodRecordApi.getPaged(params.limit, params.offset),
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

  protected openEntry(entry: MoodRecord): void {
    void this.#router.navigate(["/mood-records", entry.id]);
  }

  protected createEntry(): void {
    void this.#router.navigate(["/mood-records/create"]);
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
