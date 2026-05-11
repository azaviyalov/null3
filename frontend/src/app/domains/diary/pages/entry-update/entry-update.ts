import { Component, computed, inject, signal } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { map } from "rxjs";
import { toWritableStateSignal } from "../../../../core/utils/signal-helpers";
import { stateError } from "../../../../core/utils/state";
import { DiaryEntry, EditDiaryEntryRequest } from "../../models/entry";
import { DiaryEntryApi } from "../../services/entry-api";
import { EntryForm } from "../../components/entry-form/entry-form";

@Component({
  selector: "app-diary-entry-update",
  standalone: true,
  imports: [EntryForm],
  templateUrl: "./entry-update.html",
  styleUrl: "./entry-update.scss",
})
export class EntryUpdate {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly entryApi = inject(DiaryEntryApi);

  readonly isSubmitting = signal(false);

  private readonly entryState = toWritableStateSignal({
    trigger: this.route.params.pipe(map((params) => Number(params["id"]))),
    project: (id) => this.entryApi.getById(id),
  });

  readonly entry = computed(() => this.entryState().value);
  readonly isLoading = computed(() => this.entryState().isLoading);
  readonly errorMessage = computed(() => this.entryState().error);

  submit(payload: EditDiaryEntryRequest): void {
    const entry = this.entry();
    if (!entry) {
      return;
    }

    this.isSubmitting.set(true);

    this.entryApi.update(entry.id, payload).subscribe({
      next: (updatedEntry: DiaryEntry) => this.handleSuccess(updatedEntry),
      error: (err) => this.handleError(err),
    });
  }

  private handleSuccess(entry: DiaryEntry): void {
    this.isSubmitting.set(false);
    this.router.navigate(["/diary/entries", entry.id]);
  }

  private handleError(err: unknown): void {
    this.isSubmitting.set(false);
    this.entryState.set(stateError(err));
    console.error(err);
  }
}
