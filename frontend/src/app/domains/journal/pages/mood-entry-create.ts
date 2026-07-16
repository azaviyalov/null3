import { Component, inject, signal } from "@angular/core";
import { EditMoodEntryRequest, MoodEntry } from "../models/mood-entry";
import { MoodEntryApi } from "../services/mood-entry-api";
import { Router } from "@angular/router";
import { MoodEntryForm } from "../components/mood-entry-form";

@Component({
  selector: "app-mood-entry-create",
  standalone: true,
  imports: [MoodEntryForm],
  templateUrl: "./mood-entry-create.html",
  styleUrl: "./mood-entry-create.scss",
})
export class MoodEntryCreate {
  private readonly router = inject(Router);
  private readonly moodEntryApi = inject(MoodEntryApi);

  readonly isSubmitting = signal(false);
  readonly errorMessage = signal<string | null>(null);

  submit(payload: EditMoodEntryRequest): void {
    this.isSubmitting.set(true);
    this.errorMessage.set(null);

    this.moodEntryApi.create(payload).subscribe({
      next: (entry: MoodEntry) => this.handleSuccess(entry),
      error: (err) => this.handleError(err),
    });
  }

  private handleSuccess(entry: MoodEntry): void {
    this.isSubmitting.set(false);
    this.router.navigate(["/mood/entries", entry.id]);
  }

  private handleError(err: unknown): void {
    this.isSubmitting.set(false);
    this.errorMessage.set("Could not create the mood entry.");
    console.error(err);
  }
}
