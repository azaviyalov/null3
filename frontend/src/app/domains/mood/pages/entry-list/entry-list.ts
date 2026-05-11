import { Component, computed, inject, signal } from "@angular/core";
import { Entry } from "../../models/entry";
import { EntryApi } from "../../services/entry-api";
import { ActivatedRoute, Router } from "@angular/router";
import { EntryHistory } from "../../components/entry-history/entry-history";
import {
  toWritableSignal,
  toWritableStateSignal,
} from "../../../../core/utils/signal-helpers";
import { toObservable } from "@angular/core/rxjs-interop";
import { combineLatest, map } from "rxjs";

@Component({
  selector: "app-entry-list",
  standalone: true,
  imports: [EntryHistory],
  templateUrl: "./entry-list.html",
  styleUrl: "./entry-list.scss",
})
export class EntryList {
  static readonly defaultEntryCount = 10;

  readonly defaultEntryCount = EntryList.defaultEntryCount;

  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly entryApi = inject(EntryApi);

  readonly pageSize = signal(EntryList.defaultEntryCount);
  readonly pageOffset = signal(0);
  readonly deletedSwitch = toWritableSignal({
    trigger: this.route.queryParams.pipe(
      map((params) => params["deleted"] === "true"),
    ),
    initialValue: false,
  });

  private readonly pageState = toWritableStateSignal({
    trigger: combineLatest([
      toObservable(this.pageSize),
      toObservable(this.pageOffset),
      toObservable(this.deletedSwitch),
    ]),
    project: ([size, offset, deleted]) =>
      this.entryApi.getPaged(size, offset, deleted),
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

  openEntry(entry: Entry): void {
    this.router.navigate(["/mood/entries", entry.id]);
  }

  createEntry(): void {
    this.router.navigate(["/mood/entries/create"]);
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
    this.deletedSwitch.set(state === "deleted");
    this.pageOffset.set(0);
    this.pageSize.set(EntryList.defaultEntryCount);
    this.router.navigate([], {
      queryParams: { deleted: this.deletedSwitch() },
      queryParamsHandling: "merge",
    });
  }
}
