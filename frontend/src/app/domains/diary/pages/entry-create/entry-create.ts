import { Component, inject, signal } from "@angular/core";
import { Router } from "@angular/router";
import { DiaryEntryApi } from "../../services/entry-api";
import { DiaryEntry, EditDiaryEntryRequest } from "../../models/entry";
import { EntryForm } from "../../components/entry-form/entry-form";

@Component({
  selector: "app-diary-entry-create",
  standalone: true,
  imports: [EntryForm],
  templateUrl: "./entry-create.html",
  styleUrl: "./entry-create.scss",
})
export class EntryCreate {
  private readonly router = inject(Router);
  private readonly entryApi = inject(DiaryEntryApi);

  readonly isSubmitting = signal(false);
  readonly errorMessage = signal<string | null>(null);

  submit(payload: EditDiaryEntryRequest): void {
    this.isSubmitting.set(true);
    this.errorMessage.set(null);

    this.entryApi.create(payload).subscribe({
      next: (entry: DiaryEntry) => this.handleSuccess(entry),
      error: (err) => this.handleError(err),
    });
  }

  private handleSuccess(entry: DiaryEntry): void {
    this.isSubmitting.set(false);
    this.router.navigate(["/diary/entries", entry.id]);
  }

  private handleError(err: unknown): void {
    this.isSubmitting.set(false);
    this.errorMessage.set("Failed to create Diary Entry.");
    console.error(err);
  }
}
