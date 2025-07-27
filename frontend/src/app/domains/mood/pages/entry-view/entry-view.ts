import { Component, inject, signal, OnInit } from "@angular/core";
import { MatCardModule } from "@angular/material/card";
import { ActivatedRoute, Router } from "@angular/router";
import { Entry } from "../../models/entry";
import { EntryApi } from "../../services/entry-api";

import { MatIconModule } from "@angular/material/icon";
import { MatButtonModule } from "@angular/material/button";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { EntryCard } from "../../components/entry-card/entry-card";

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
export class EntryView implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly entryApi = inject(EntryApi);

  readonly entry = signal<Entry | null>(null);
  readonly isLoading = signal(true);

  ngOnInit(): void {
    this.loadEntry(Number(this.route.snapshot.paramMap.get("id")));
  }

  readonly errorMessage = signal<string | null>(null);

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

  editEntry(): void {
    const entry = this.entry();
    if (!entry || entry.deletedAt) {
      this.errorMessage.set("Cannot edit an already deleted entry");
      return;
    }
    this.router.navigate(["/mood/entries", entry.id, "update"]);
  }

  deleteEntry(): void {
    const entry = this.entry();
    if (!entry || entry.deletedAt) {
      this.errorMessage.set("Cannot delete an already deleted entry");
      return;
    }
    this.entryApi.delete(entry.id).subscribe({
      next: (deletedEntry) => {
        this.entry.set(deletedEntry);
        this.errorMessage.set(null);
      },
      error: (err) => {
        this.errorMessage.set(err?.message ?? "Failed to delete entry");
      },
    });
  }

  restoreEntry(): void {
    const entry = this.entry();
    if (!entry?.deletedAt) {
      this.errorMessage.set("Cannot restore an active entry");
      return;
    }
    this.entryApi.restore(entry.id).subscribe({
      next: (restoredEntry) => {
        this.entry.set(restoredEntry);
        this.errorMessage.set(null);
      },
      error: (err) => {
        this.errorMessage.set(err?.message ?? "Failed to restore entry");
      },
    });
  }

  goBack(): void {
    history.back();
  }
}
