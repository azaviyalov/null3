import { Component, inject } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatIconModule } from "@angular/material/icon";
import { RouterModule } from "@angular/router";
import { EntryCardGrid } from "../../../domains/mood/components/entry-card-grid/entry-card-grid";
import { toSignal } from "@angular/core/rxjs-interop";
import { EntryApi } from "../../../domains/mood/services/entry-api";

@Component({
  selector: "app-home",
  imports: [
    RouterModule,
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    EntryCardGrid,
  ],
  templateUrl: "./home.html",
  styleUrl: "./home.scss",
})
export class Home {
  private readonly entryApi = inject(EntryApi);
  private readonly limit = 4;
  readonly entriesPage = toSignal(this.entryApi.getPaged(this.limit), {
    initialValue: null,
  });
}
