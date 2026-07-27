import { Component, computed, inject, signal } from "@angular/core";
import { Router } from "@angular/router";
import { toObservable } from "@angular/core/rxjs-interop";
import { toWritableStateSignal } from "../../../core/utils/signal-helpers";
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

  private readonly router = inject(Router);
  private readonly diaryEntryApi = inject(DiaryEntryApi);

  readonly pageOffset = signal(0);

  private readonly pageState = toWritableStateSignal({
    trigger: toObservable(this.pageOffset),
    project: (offset) =>
      this.diaryEntryApi.getPaged(this.defaultPageSize, offset),
  });

  readonly page = computed(() => this.pageState().value);
  readonly isLoading = computed(() => this.pageState().isLoading);
  readonly totalCount = computed(() => this.page()?.totalCount ?? 0);
  readonly canGoPrevious = computed(() => this.pageOffset() > 0);
  readonly canGoNext = computed(
    () => this.pageOffset() + this.defaultPageSize < this.totalCount(),
  );

  openEntry(entry: DiaryEntry): void {
    this.router.navigate(["/diary-entries", entry.id]);
  }

  createEntry(): void {
    this.router.navigate(["/diary-entries/create"]);
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
