import { Component, computed, inject } from "@angular/core";
import { MatCardModule } from "@angular/material/card";
import { ActivatedRoute, Router } from "@angular/router";
import { EntryApi } from "../../services/entry-api";
import { MatIconModule } from "@angular/material/icon";
import { MatButtonModule } from "@angular/material/button";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { EntryCard } from "../../components/entry-card/entry-card";
import { map } from "rxjs";
import { toWritableStateSignal } from "../../../../core/utils/signal-helpers";
import { stateError, stateSuccess } from "../../../../core/utils/state";
@Component({
  selector: "app-entry-view",
  standalone: true,
  imports: [
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    MatProgressSpinnerModule,
    EntryCard,
  ],
  templateUrl: "./entry-view.html",
  styleUrl: "./entry-view.scss",
})
export class EntryView {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly entryApi = inject(EntryApi);

  private readonly entryState = toWritableStateSignal({
    trigger: this.route.params.pipe(map((params) => Number(params["id"]))),
    project: (id) => this.entryApi.getById(id),
  });

  readonly entry = computed(() => this.entryState().value);
  readonly isLoading = computed(() => this.entryState().isLoading);
  readonly errorMessage = computed(() => this.entryState().error);

  editEntry(): void {
    const entry = this.entry();
    if (!entry || entry.deletedAt) {
      this.entryState.set(stateError("Cannot edit a deleted entry"));
      return;
    }

    this.router.navigate(["/mood/entries", entry.id, "update"]);
  }

  deleteEntry(): void {
    const entry = this.entry();
    if (!entry || entry.deletedAt) {
      this.entryState.set(stateError("Cannot delete a deleted entry"));
      return;
    }

    this.entryApi.delete(entry.id).subscribe({
      next: (deletedEntry) =>
        this.entryState.set(stateSuccess(deletedEntry)),
      error: (err) => this.entryState.set(stateError(err)),
    });
  }

  restoreEntry(): void {
    const entry = this.entry();
    if (!entry?.deletedAt) {
      this.entryState.set(stateError("Cannot restore an active entry"));
      return;
    }

    this.entryApi.restore(entry.id).subscribe({
      next: (restoredEntry) => {
        this.entryState.set(stateSuccess(restoredEntry));
      },
      error: (err) => {
        this.entryState.set(stateError(err));
      },
    });
  }

  goBack(): void {
    history.back();
  }
}
