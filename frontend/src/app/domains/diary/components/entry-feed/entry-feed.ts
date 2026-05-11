import { CommonModule } from "@angular/common";
import { Component, computed, input, output } from "@angular/core";
import { DiaryEntry } from "../../models/entry";

interface EntryGroup {
  readonly key: string;
  readonly label: string;
  readonly entries: DiaryEntry[];
}

@Component({
  selector: "app-diary-entry-feed",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./entry-feed.html",
  styleUrl: "./entry-feed.scss",
})
export class EntryFeed {
  readonly skeleton = input(false);
  readonly skeletonCount = input(10);
  readonly entries = input<DiaryEntry[] | null>(null);
  readonly showOpen = input(false);
  readonly showEdit = input(false);
  readonly showDelete = input(false);
  readonly showRestore = input(false);
  readonly emptyMessage = input("No Diary Entries found");

  readonly groupedEntries = computed<EntryGroup[]>(() =>
    groupEntriesByDate(this.entries() ?? []),
  );
  readonly skeletonRows = computed(() =>
    Array.from({ length: this.skeletonCount() }, (_, index) => index),
  );
  readonly showActions = computed(
    () =>
      this.showOpen() ||
      this.showEdit() ||
      this.showDelete() ||
      this.showRestore(),
  );

  readonly open = output<DiaryEntry>();
  readonly edit = output<DiaryEntry>();
  readonly delete = output<DiaryEntry>();
  readonly restore = output<DiaryEntry>();

  trackByGroup = (_: number, group: EntryGroup): string => group.key;
  trackByEntry = (_: number, entry: DiaryEntry): number => entry.id;
}

function groupEntriesByDate(entries: DiaryEntry[]): EntryGroup[] {
  const groups = new Map<string, EntryGroup>();

  for (const entry of entries) {
    const key = dateKey(entry.occurredAt);
    const existingGroup = groups.get(key);
    if (existingGroup) {
      existingGroup.entries.push(entry);
      continue;
    }

    groups.set(key, {
      key,
      label: formatGroupLabel(entry.occurredAt),
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
