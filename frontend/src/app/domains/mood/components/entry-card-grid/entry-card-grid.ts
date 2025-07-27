import { Component, input, computed, output } from "@angular/core";
import { CommonModule } from "@angular/common";
import { EntryCard } from "../entry-card/entry-card";
import { Entry } from "../../models/entry";

@Component({
  selector: "app-entry-card-grid",
  standalone: true,
  imports: [CommonModule, EntryCard],
  templateUrl: "./entry-card-grid.html",
  styleUrl: "./entry-card-grid.scss",
})
export class EntryCardGrid {
  readonly skeleton = input(false);
  readonly skeletonCount = input(10);
  readonly entries = input<Entry[] | null>(null);
  readonly entriesOrSkeleton = computed<(Entry | null)[] | null>(() =>
    this.skeleton()
      ? Array.from({ length: this.skeletonCount() }, () => null)
      : this.entries(),
  );
  readonly columns = input(1);
  readonly showOpen = input(false);
  readonly showEdit = input(false);
  readonly showDelete = input(false);
  readonly showRestore = input(false);

  readonly gridStyle = computed(() => ({
    "grid-template-columns": `repeat(${this.columns()}, 1fr)`,
  }));

  readonly open = output<Entry>();
  readonly edit = output<Entry>();
  readonly delete = output<Entry>();
  readonly restore = output<Entry>();

  trackByEntryOrIndex = (_: number, entry: Entry | null) =>
    entry ? entry.id : _;
}
