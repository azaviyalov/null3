import { Component, computed, inject, signal } from "@angular/core";
import { MoodRecord } from "../models/mood-record";
import { MoodRecordApi } from "../services/mood-record-api";
import { ActivatedRoute, Router } from "@angular/router";
import { MoodRecordHistory } from "../components/mood-record-history";
import {
  toWritableSignal,
  toWritableStateSignal,
} from "../../../core/utils/signal-helpers";
import { toObservable } from "@angular/core/rxjs-interop";
import { combineLatest, map } from "rxjs";

@Component({
  selector: "app-mood-record-list",
  standalone: true,
  imports: [MoodRecordHistory],
  templateUrl: "./mood-record-list.html",
  styleUrl: "./mood-record-list.scss",
})
export class MoodRecordList {
  readonly defaultPageSize = 10;

  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly moodRecordApi = inject(MoodRecordApi);

  readonly pageSize = signal(this.defaultPageSize);
  readonly pageOffset = signal(0);
  readonly showDeleted = toWritableSignal({
    trigger: this.route.queryParams.pipe(
      map((params) => params["deleted"] === "true"),
    ),
    initialValue: false,
  });

  private readonly pageState = toWritableStateSignal({
    trigger: combineLatest([
      toObservable(this.pageSize),
      toObservable(this.pageOffset),
      toObservable(this.showDeleted),
    ]),
    project: ([size, offset, deleted]) =>
      this.moodRecordApi.getPaged(size, offset, deleted),
  });

  readonly page = computed(() => this.pageState().value);
  readonly isLoading = computed(() => this.pageState().isLoading);
  readonly pageSizeOptions = [5, 10, 25, 100];
  readonly totalCount = computed(() => this.page()?.totalCount ?? 0);
  readonly currentPage = computed(() =>
    this.totalCount() ? Math.floor(this.pageOffset() / this.pageSize()) + 1 : 0,
  );
  readonly totalPages = computed(() =>
    this.totalCount() ? Math.ceil(this.totalCount() / this.pageSize()) : 0,
  );
  readonly pageStart = computed(() =>
    this.totalCount() ? this.pageOffset() + 1 : 0,
  );
  readonly pageEnd = computed(() =>
    this.totalCount()
      ? Math.min(this.pageOffset() + this.pageSize(), this.totalCount())
      : 0,
  );
  readonly canGoPrevious = computed(() => this.pageOffset() > 0);
  readonly canGoNext = computed(() => this.pageEnd() < this.totalCount());

  openEntry(entry: MoodRecord): void {
    this.router.navigate(["/mood/records", entry.id]);
  }

  createEntry(): void {
    this.router.navigate(["/mood/records/create"]);
  }

  changePageSize(nextPageSize: string): void {
    const parsedPageSize = Number(nextPageSize);
    if (!this.pageSizeOptions.includes(parsedPageSize)) {
      return;
    }

    this.pageSize.set(parsedPageSize);
    this.pageOffset.set(0);
  }

  goToPreviousPage(): void {
    if (!this.canGoPrevious()) {
      return;
    }

    this.pageOffset.update((value) => Math.max(0, value - this.pageSize()));
  }

  goToNextPage(): void {
    if (!this.canGoNext()) {
      return;
    }

    this.pageOffset.update((value) => value + this.pageSize());
  }

  setDeletedState(state: "active" | "deleted"): void {
    this.showDeleted.set(state === "deleted");
    this.pageOffset.set(0);
    this.pageSize.set(this.defaultPageSize);
    this.router.navigate([], {
      queryParams: { deleted: this.showDeleted() },
      queryParamsHandling: "merge",
    });
  }
}
