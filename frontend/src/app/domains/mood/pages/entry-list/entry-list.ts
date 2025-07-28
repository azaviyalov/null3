import { Component, inject, signal, effect, OnInit } from "@angular/core";
import { MatCardModule } from "@angular/material/card";
import { Entry, PaginatedEntryList } from "../../models/entry";
import { EntryApi } from "../../services/entry-api";
import { MatIconModule } from "@angular/material/icon";
import { ActivatedRoute, Router } from "@angular/router";
import { MatButtonModule } from "@angular/material/button";
import { PageEvent, MatPaginatorModule } from "@angular/material/paginator";
import { MatButtonToggleModule } from "@angular/material/button-toggle";
import { EntryCardGrid } from "../../components/entry-card-grid/entry-card-grid";

@Component({
  selector: "app-entry-list",
  standalone: true,
  imports: [
    MatButtonModule,
    MatButtonToggleModule,
    MatCardModule,
    MatIconModule,
    MatPaginatorModule,
    EntryCardGrid
],
  templateUrl: "./entry-list.html",
  styleUrl: "./entry-list.scss",
})
export class EntryList implements OnInit {
  readonly defaultCardCount = 10;

  readonly pageSize = signal<number>(10);
  readonly pageOffset = signal<number>(0);
  readonly deletedSwitch = signal<boolean>(false);
  readonly entriesPage = signal<PaginatedEntryList | null>(null);
  readonly isLoading = signal(true);

  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly api = inject(EntryApi);

  constructor() {
    effect(() => {
      this.isLoading.set(true);
      const size = this.pageSize();
      const offset = this.pageOffset();
      const deleted = this.deletedSwitch();
      this.api.getPaged(size, offset, deleted).subscribe((result) => {
        this.entriesPage.set(result);
        this.isLoading.set(false);
      });
    });
  }

  ngOnInit(): void {
    this.activatedRoute.queryParams.subscribe((params) => {
      const deleted = params["deleted"];
      if (deleted !== undefined) {
        this.deletedSwitch.set(deleted === "true");
      }
    });
  }

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

  setDeletedState(state: 'active' | 'deleted'): void {
    this.deletedSwitch.set(state === "deleted");
    this.pageOffset.set(0);
    this.pageSize.set(10);
    this.router.navigate([], {
      queryParams: { deleted: this.deletedSwitch() },
      queryParamsHandling: "merge",
    });
  }
}
