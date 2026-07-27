import { groupEntriesByDate } from "./date-grouping";

interface Entry {
  readonly id: number;
  readonly date: Date;
}

describe("groupEntriesByDate", () => {
  const today = new Date(2024, 0, 15, 12);

  it("groups entries by the selected local date and preserves their order", () => {
    const entries: Entry[] = [
      { id: 1, date: new Date(2024, 0, 15, 8) },
      { id: 2, date: new Date(2024, 0, 14, 20) },
      { id: 3, date: new Date(2024, 0, 15, 18) },
    ];

    const groups = groupEntriesByDate(entries, (entry) => entry.date, today);

    expect(groups).toEqual([
      {
        key: "2024-01-15",
        label: "Today",
        entries: [entries[0], entries[2]],
      },
      {
        key: "2024-01-14",
        label: "Yesterday",
        entries: [entries[1]],
      },
    ]);
  });

  it("formats recent, same-year, and earlier-year dates", () => {
    const recent = new Date(2024, 0, 10, 12);
    const sameYear = new Date(2024, 0, 2, 12);
    const earlierYear = new Date(2023, 11, 30, 12);

    const groups = groupEntriesByDate(
      [
        { id: 1, date: recent },
        { id: 2, date: sameYear },
        { id: 3, date: earlierYear },
      ],
      (entry) => entry.date,
      today,
    );

    expect(groups.map((group) => group.label)).toEqual([
      new Intl.DateTimeFormat(undefined, {
        weekday: "long"
      }).format(recent),
      new Intl.DateTimeFormat(undefined, {
        month: "long",
        day: "numeric",
      }).format(sameYear),
      new Intl.DateTimeFormat(undefined, {
        month: "long",
        day: "numeric",
        year: "numeric",
      }).format(earlierYear),
    ]);
  });
});
