export interface DateGroup<T> {
  readonly key: string;
  readonly label: string;
  readonly entries: T[];
}

export function groupEntriesByDate<T>(
  entries: readonly T[],
  dateOf: (entry: T) => Date,
  today = new Date(),
): DateGroup<T>[] {
  const groups = new Map<string, DateGroup<T>>();

  for (const entry of entries) {
    const date = dateOf(entry);
    const key = dateKey(date);
    const existingGroup = groups.get(key);

    if (existingGroup) {
      existingGroup.entries.push(entry);
      continue;
    }

    groups.set(key, {
      key,
      label: formatGroupLabel(date, today),
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

function formatGroupLabel(date: Date, today: Date): string {
  const startOfToday = startOfDay(today);
  const target = startOfDay(date);
  const diffInDays = Math.round(
    (startOfToday.getTime() - target.getTime()) / 86400000,
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

  const sameYear = startOfToday.getFullYear() === target.getFullYear();
  return new Intl.DateTimeFormat(undefined, {
    month: "long",
    day: "numeric",
    ...(sameYear ? {} : { year: "numeric" }),
  }).format(date);
}

function startOfDay(date: Date): Date {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate());
}
