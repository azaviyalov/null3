import { Component, inject, signal } from "@angular/core";
import { EditMoodRecordRequest, MoodRecord } from "../models/mood-record";
import { MoodRecordApi } from "../services/mood-record-api";
import { Router } from "@angular/router";
import { MoodRecordForm } from "../components/mood-record-form";

@Component({
  selector: "app-mood-record-create",
  standalone: true,
  imports: [MoodRecordForm],
  templateUrl: "./mood-record-create.html",
  styleUrl: "./mood-record-create.scss",
})
export class MoodRecordCreate {
  private readonly router = inject(Router);
  private readonly moodRecordApi = inject(MoodRecordApi);

  readonly isSubmitting = signal(false);
  readonly errorMessage = signal<string | null>(null);

  submit(payload: EditMoodRecordRequest): void {
    this.isSubmitting.set(true);
    this.errorMessage.set(null);

    this.moodRecordApi.create(payload).subscribe({
      next: (entry: MoodRecord) => this.handleSuccess(entry),
      error: (err) => this.handleError(err),
    });
  }

  private handleSuccess(entry: MoodRecord): void {
    this.isSubmitting.set(false);
    this.router.navigate(["/mood-records", entry.id]);
  }

  private handleError(err: unknown): void {
    this.isSubmitting.set(false);
    this.errorMessage.set("Could not create the mood record.");
    console.error(err);
  }
}
