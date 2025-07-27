import { Component, inject, signal } from "@angular/core";
import { EntryApi } from "../../services/entry-api";
import { EditEntryRequest, Entry } from "../../models/entry";
import { Router } from "@angular/router";
import { EntryForm } from "../../components/entry-form/entry-form";

@Component({
  selector: "app-entry-create",
  standalone: true,
  imports: [EntryForm],
  templateUrl: "./entry-create.html",
  styleUrl: "./entry-create.scss",
})
export class EntryCreate {
  private readonly router = inject(Router);
  private readonly entryApi = inject(EntryApi);

  readonly isSubmitting = signal(false);
  readonly errorMessage = signal<string | null>(null);

  submit(payload: EditEntryRequest): void {
    this.isSubmitting.set(true);
    this.errorMessage.set(null);

    this.entryApi.create(payload).subscribe({
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
    this.errorMessage.set("Failed to create entry.");
    console.error(err);
  }
}
