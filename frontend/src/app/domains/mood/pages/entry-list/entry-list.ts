import { Component, inject, signal, effect } from "@angular/core";
import { MatCardModule } from "@angular/material/card";
import { Entry, PaginatedEntryList } from "../../models/entry";
import { EntryApi } from "../../services/entry-api";
import { MatIconModule } from "@angular/material/icon";
import { Router } from "@angular/router";
import { MatButtonModule } from "@angular/material/button";
import { PageEvent } from "@angular/material/paginator";
import { EntryCardGridPaginated } from "../../components/entry-card-grid-paginated/entry-card-grid-paginated";

@Component({
  selector: "app-entry-list",
  standalone: true,
  imports: [
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    EntryCardGridPaginated,
  ],
  templateUrl: "./entry-list.html",
  styleUrl: "./entry-list.scss",
})
export class EntryList {
  readonly pageSize = signal<number>(10);
  readonly pageOffset = signal<number>(0);
  readonly entriesPage = signal<PaginatedEntryList | null>(null);

  private readonly router = inject(Router);
  private readonly api = inject(EntryApi);

  isLoading = signal(true);

  constructor() {
    effect(() => {
      this.isLoading.set(true);
      const size = this.pageSize();
      const offset = this.pageOffset();
      this.api.getPaged(size, offset).subscribe((result) => {
        this.entriesPage.set(result);
        this.isLoading.set(false);
      });
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
}
