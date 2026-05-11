import { Component, computed, input, output } from "@angular/core";
import { CommonModule } from "@angular/common";
import { Entry } from "../../models/entry";
import { feelingLabel } from "../../utils/entry-presenter";

@Component({
  selector: "app-entry-detail",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./entry-detail.html",
  styleUrl: "./entry-detail.scss",
})
export class EntryDetail {
  readonly skeleton = input(false);
  readonly entry = input<Entry | null>(null);
  readonly showOpen = input(false);
  readonly showEdit = input(false);
  readonly showDelete = input(false);
  readonly showRestore = input(false);

  readonly showActions = computed(
    () =>
      this.showOpen() ||
      this.showEdit() ||
      this.showDelete() ||
      this.showRestore(),
  );

  readonly showFooter = computed(
    () =>
      this.entry() &&
      (this.entry()?.deletedAt ||
        this.entry()?.createdAt?.getTime() !==
          this.entry()?.updatedAt?.getTime()) &&
      !this.skeleton(),
  );

  readonly updated = computed(
    () =>
      this.entry() &&
      this.entry()?.updatedAt.getTime() !== this.entry()?.createdAt.getTime(),
  );
  readonly deleted = computed(() => !!this.entry()?.deletedAt);
  readonly feelingLabel = feelingLabel;

  readonly open = output<void>();
  readonly edit = output<void>();
  readonly delete = output<void>();
  readonly restore = output<void>();
}
