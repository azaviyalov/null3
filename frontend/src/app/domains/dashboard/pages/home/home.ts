import { Component, inject } from "@angular/core";
import { Router, RouterModule } from "@angular/router";
import { toSignal } from "@angular/core/rxjs-interop";
import { EntryHistory } from "../../../mood/components/entry-history/entry-history";
import { EntryApi } from "../../../mood/services/entry-api";
import { Entry } from "../../../mood/models/entry";

@Component({
  selector: "app-home",
  imports: [RouterModule, EntryHistory],
  templateUrl: "./home.html",
  styleUrl: "./home.scss",
})
export class Home {
  private readonly entryApi = inject(EntryApi);
  private readonly router = inject(Router);
  private readonly limit = 6;

  readonly entriesPage = toSignal(this.entryApi.getPaged(this.limit), {
    initialValue: null,
  });

  openEntry(entry: Entry): void {
    this.router.navigate(["/mood/entries", entry.id]);
  }
}
