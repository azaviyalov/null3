import { Component, inject, signal } from "@angular/core";
import { Router } from "@angular/router";
import { DiaryEntryApi } from "../services/diary-entry-api";
import { DiaryEntry, EditDiaryEntryRequest } from "../models/diary-entry";
import { DiaryEntryForm } from "../components/diary-entry-form";

@Component({
  selector: "app-diary-entry-create",
  standalone: true,
  imports: [DiaryEntryForm],
  templateUrl: "./diary-entry-create.html",
  styleUrl: "./diary-entry-create.scss",
})
export class DiaryEntryCreate {
  private readonly router = inject(Router);
  private readonly diaryEntryApi = inject(DiaryEntryApi);

  readonly isSubmitting = signal(false);
  readonly errorMessage = signal<string | null>(null);

  submit(payload: EditDiaryEntryRequest): void {
    this.isSubmitting.set(true);
    this.errorMessage.set(null);

    this.diaryEntryApi.create(payload).subscribe({
      next: (entry: DiaryEntry) => this.handleSuccess(entry),
      error: (err) => this.handleError(err),
    });
  }

  private handleSuccess(entry: DiaryEntry): void {
    this.isSubmitting.set(false);
    this.router.navigate(["/diary-entries", entry.id]);
  }

  private handleError(err: unknown): void {
    this.isSubmitting.set(false);
    this.errorMessage.set("Could not create the diary entry.");
    console.error(err);
  }
}
