import { Component, computed, inject, signal } from "@angular/core";
import { MoodRecord } from "../models/mood-record";
import { MoodRecordApi } from "../services/mood-record-api";
import { Router } from "@angular/router";
import { MoodRecordHistory } from "../components/mood-record-history";
import { toWritableStateSignal } from "../../../core/utils/signal-helpers";
import { toObservable } from "@angular/core/rxjs-interop";

@Component({
  selector: "app-mood-record-list",
  standalone: true,
  imports: [MoodRecordHistory],
  templateUrl: "./mood-record-list.html",
  styleUrl: "./mood-record-list.scss",
})
export class MoodRecordList {
  readonly defaultPageSize = 10;

  private readonly router = inject(Router);
  private readonly moodRecordApi = inject(MoodRecordApi);

  readonly pageOffset = signal(0);

  private readonly pageState = toWritableStateSignal({
    trigger: toObservable(this.pageOffset),
    project: (offset) =>
      this.moodRecordApi.getPaged(this.defaultPageSize, offset),
  });

  readonly page = computed(() => this.pageState().value);
  readonly isLoading = computed(() => this.pageState().isLoading);
  readonly totalCount = computed(() => this.page()?.totalCount ?? 0);
  readonly canGoPrevious = computed(() => this.pageOffset() > 0);
  readonly canGoNext = computed(
    () => this.pageOffset() + this.defaultPageSize < this.totalCount(),
  );

  openEntry(entry: MoodRecord): void {
    this.router.navigate(["/mood-records", entry.id]);
  }

  createEntry(): void {
    this.router.navigate(["/mood-records/create"]);
  }

  goToPreviousPage(): void {
    if (!this.canGoPrevious()) {
      return;
    }

    this.pageOffset.update((value) =>
      Math.max(0, value - this.defaultPageSize),
    );
  }

  goToNextPage(): void {
    if (!this.canGoNext()) {
      return;
    }

    this.pageOffset.update((value) => value + this.defaultPageSize);
  }
}
