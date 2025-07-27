import { Component, input, output } from "@angular/core";
import { Entry, PaginatedEntryList } from "../../models/entry";
import { MatPaginatorModule, PageEvent } from "@angular/material/paginator";
import { EntryCardGrid } from "../entry-card-grid/entry-card-grid";

@Component({
  selector: "app-entry-card-grid-paginated",
  imports: [MatPaginatorModule, EntryCardGrid],
  templateUrl: "./entry-card-grid-paginated.html",
  styleUrl: "./entry-card-grid-paginated.scss",
})
export class EntryCardGridPaginated {
  readonly defaultCardCount = 10;

  readonly entriesPage = input<PaginatedEntryList | null>(null);
  readonly columns = input(1);
  readonly showOpen = input(false);
  readonly showEdit = input(false);
  readonly showDelete = input(false);
  readonly showRestore = input(false);
  readonly pageSize = input<number>(this.defaultCardCount);
  readonly skeleton = input(true);

  readonly changePage = output<PageEvent>();
  readonly open = output<Entry>();
  readonly edit = output<Entry>();
  readonly delete = output<Entry>();
  readonly restore = output<Entry>();
}
