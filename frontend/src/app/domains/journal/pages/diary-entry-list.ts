import { Component, computed, inject, signal } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { combineLatest, map } from "rxjs";
import { toObservable } from "@angular/core/rxjs-interop";
import {
  toWritableSignal,
  toWritableStateSignal,
} from "../../../core/utils/signal-helpers";
import { DiaryEntryApi } from "../services/diary-entry-api";
import { DiaryEntry } from "../models/diary-entry";
import { DiaryEntryFeed } from "../components/diary-entry-feed";

@Component({
  selector: "app-diary-entry-list",
  standalone: true,
  imports: [DiaryEntryFeed],
  templateUrl: "./diary-entry-list.html",
  styleUrl: "./diary-entry-list.scss",
})
export class DiaryEntryList {
  readonly defaultPageSize = 10;

  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly diaryEntryApi = inject(DiaryEntryApi);

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
      this.diaryEntryApi.getPaged(size, offset, deleted),
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

  openEntry(entry: DiaryEntry): void {
    this.router.navigate(["/diary-entries", entry.id]);
  }

  createEntry(): void {
    this.router.navigate(["/diary-entries/create"]);
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
