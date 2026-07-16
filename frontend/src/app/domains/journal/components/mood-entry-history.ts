import { Component, computed, input, output } from "@angular/core";
import { CommonModule } from "@angular/common";
import { MoodEntry } from "../models/mood-entry";
import { feelingLabel } from "../utils/feeling-presenter";

interface EntryGroup {
  readonly key: string;
  readonly label: string;
  readonly entries: MoodEntry[];
}

@Component({
  selector: "app-mood-entry-history",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./mood-entry-history.html",
  styleUrl: "./mood-entry-history.scss",
})
export class MoodEntryHistory {
  readonly skeleton = input(false);
  readonly skeletonCount = input(10);
  readonly entries = input<MoodEntry[] | null>(null);
  readonly showOpen = input(false);
  readonly showEdit = input(false);
  readonly showDelete = input(false);
  readonly showRestore = input(false);
  readonly emptyMessage = input("No mood entries found.");

  readonly groupedEntries = computed<EntryGroup[]>(() =>
    groupEntriesByDate(this.entries() ?? []),
  );
  readonly skeletonRows = computed(() =>
    Array.from({ length: this.skeletonCount() }, (_, index) => index),
  );

  readonly open = output<MoodEntry>();
  readonly edit = output<MoodEntry>();
  readonly delete = output<MoodEntry>();
  readonly restore = output<MoodEntry>();

  readonly showActions = computed(
    () =>
      this.showOpen() ||
      this.showEdit() ||
      this.showDelete() ||
      this.showRestore(),
  );
  readonly feelingLabel = feelingLabel;

  trackByGroup = (_: number, group: EntryGroup): string => group.key;
  trackByEntry = (_: number, entry: MoodEntry): number => entry.id;
}

function groupEntriesByDate(entries: MoodEntry[]): EntryGroup[] {
  const groups = new Map<string, EntryGroup>();

  for (const entry of entries) {
    const key = dateKey(entry.createdAt);
    const existingGroup = groups.get(key);

    if (existingGroup) {
      existingGroup.entries.push(entry);
      continue;
    }

    groups.set(key, {
      key,
      label: formatGroupLabel(entry.createdAt),
      entries: [entry],
    });
  }

  return Array.from(groups.values());
}

function dateKey(date: Date): string {
  const year = date.getFullYear();
  const month = `${date.getMonth() + 1}`.padStart(2, "0");
  const day = `${date.getDate()}`.padStart(2, "0");
  return `${year}-${month}-${day}`;
}

function formatGroupLabel(date: Date): string {
  const today = startOfDay(new Date());
  const target = startOfDay(date);
  const diffInDays = Math.round(
    (today.getTime() - target.getTime()) / 86400000,
  );

  if (diffInDays === 0) {
    return "Today";
  }

  if (diffInDays === 1) {
    return "Yesterday";
  }

  if (diffInDays >= 0 && diffInDays < 7) {
    return new Intl.DateTimeFormat(undefined, {
      weekday: "long",
    }).format(date);
  }

  const sameYear = today.getFullYear() === target.getFullYear();
  return new Intl.DateTimeFormat(undefined, {
    month: "long",
    day: "numeric",
    ...(sameYear ? {} : { year: "numeric" }),
  }).format(date);
}

function startOfDay(date: Date): Date {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate());
}
