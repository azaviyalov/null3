import { Component, inject, signal, computed } from "@angular/core";
import { MatCardModule } from "@angular/material/card";
import { Entry } from "../../models/entry";
import { EntryApi } from "../../services/entry-api";
import { MatIconModule } from "@angular/material/icon";
import { Router } from "@angular/router";
import { MatButtonModule } from "@angular/material/button";
import { PageEvent, MatPaginatorModule } from "@angular/material/paginator";
import { MatButtonToggleModule } from "@angular/material/button-toggle";
import { EntryCardGrid } from "../../components/entry-card-grid/entry-card-grid";
import {
  toWritableSignal,
  toWritableStateSignal,
} from "../../../../core/utils/signal-helpers";
import { toObservable } from "@angular/core/rxjs-interop";
import { combineLatest, map } from "rxjs";

@Component({
  selector: "app-entry-list",
  standalone: true,
  imports: [
    MatButtonModule,
    MatButtonToggleModule,
    MatCardModule,
    MatIconModule,
    MatPaginatorModule,
    EntryCardGrid,
  ],
  templateUrl: "./entry-list.html",
  styleUrl: "./entry-list.scss",
})
export class EntryList {
  static readonly defaultCardCount = 10;

  readonly EntryList = EntryList;

  private readonly router = inject(Router);
  private readonly api = inject(EntryApi);

  readonly pageSize = signal(EntryList.defaultCardCount);
  readonly pageOffset = signal(0);
  readonly deletedSwitch = toWritableSignal({
    trigger: this.router.routerState.root.queryParams.pipe(
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
      this.api.getPaged(size, offset, deleted),
  });

  readonly page = computed(() => this.pageState().value);
  readonly isLoading = computed(() => this.pageState().isLoading);

  openEntry(entry: Entry): void {
    this.router.navigate(["/mood/entries", entry.id]);
  }

  createEntry(): void {
    this.router.navigate(["/mood/entries/create"]);
  }

  changePage(event: PageEvent): void {
    this.pageSize.set(event.pageSize);
    this.pageOffset.set(event.pageIndex * event.pageSize);
  }

  setDeletedState(state: "active" | "deleted"): void {
    this.deletedSwitch.set(state === "deleted");
    this.pageOffset.set(0);
    this.pageSize.set(EntryList.defaultCardCount);
    this.router.navigate([], {
      queryParams: { deleted: this.deletedSwitch() },
      queryParamsHandling: "merge",
    });
  }
}
