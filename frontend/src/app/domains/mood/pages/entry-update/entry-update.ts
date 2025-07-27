import { Component, inject, signal, OnInit } from "@angular/core";
import { EditEntryRequest, Entry } from "../../models/entry";
import { ActivatedRoute, Router } from "@angular/router";
import { EntryApi } from "../../services/entry-api";
import { EntryForm } from "../../components/entry-form/entry-form";

@Component({
  selector: "app-entry-update",
  standalone: true,
  imports: [EntryForm],
  templateUrl: "./entry-update.html",
  styleUrl: "./entry-update.scss",
})
export class EntryUpdate implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly entryApi = inject(EntryApi);

  readonly entry = signal<Entry | null>(null);
  readonly isLoading = signal(true);

  readonly isSubmitting = signal(false);
  readonly errorMessage = signal<string | null>(null);

  ngOnInit(): void {
    this.loadEntry(Number(this.route.snapshot.paramMap.get("id")));
  }

  private loadEntry(id: number): void {
    this.isLoading.set(true);
    this.entryApi.getById(id).subscribe({
      next: (entry) => {
        this.entry.set(entry);
        this.isLoading.set(false);
      },
      error: (err) => {
        this.errorMessage.set(err?.message ?? "Failed to load entry");
        this.isLoading.set(false);
      },
    });
  }

  submit(payload: EditEntryRequest): void {
    const entry = this.entry();
    if (!entry) {
      this.errorMessage.set("No entry to update");
      return;
    }

    this.isSubmitting.set(true);
    this.errorMessage.set(null);

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
    this.errorMessage.set("Failed to update entry.");
    console.error(err);
  }
}
