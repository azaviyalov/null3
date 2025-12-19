import { Component, inject, signal, computed } from "@angular/core";
import { EditEntryRequest, Entry } from "../../models/entry";
import { ActivatedRoute, Router } from "@angular/router";
import { EntryApi } from "../../services/entry-api";
import { EntryForm } from "../../components/entry-form/entry-form";
import { map } from "rxjs";
import { toWritableStateSignal } from "../../../../core/utils/signal-helpers";
import { stateError } from "../../../../core/utils/state";

@Component({
  selector: "app-entry-update",
  standalone: true,
  imports: [EntryForm],
  templateUrl: "./entry-update.html",
  styleUrl: "./entry-update.scss",
})
export class EntryUpdate {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly entryApi = inject(EntryApi);

  readonly isSubmitting = signal(false);

  private readonly entryState = toWritableStateSignal({
    trigger: this.route.params.pipe(map((params) => Number(params["id"]))),
    project: (id) => this.entryApi.getById(id),
  });

  readonly entry = computed(() => this.entryState().value);
  readonly isLoading = computed(() => this.entryState().isLoading);

  submit(payload: EditEntryRequest): void {
    const entry = this.entry();
    if (!entry) {
      return;
    }

    this.isSubmitting.set(true);

    this.entryApi.update(entry.id, payload).subscribe({
      next: (entry: Entry) => this.handleSuccess(entry),
      error: (err) => this.handleError(err),
    });
  }

  private handleSuccess(entry: Entry): void {
    this.isSubmitting.set(false);
    this.router.navigate(["/mood/entries", entry.id]);
  }

  private handleError(err: unknown): void {
    this.isSubmitting.set(false);
    this.entryState.set(stateError(err));
    console.error(err);
  }
}
