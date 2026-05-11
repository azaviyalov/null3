import { Component, inject } from "@angular/core";
import { Router, RouterModule } from "@angular/router";
import { toSignal } from "@angular/core/rxjs-interop";
import { EntryHistory } from "../../../mood/components/entry-history/entry-history";
import { EntryApi } from "../../../mood/services/entry-api";
import { Entry } from "../../../mood/models/entry";
import { DiaryEntryApi } from "../../../diary/services/entry-api";
import { DiaryEntry } from "../../../diary/models/entry";
import { EntryFeed } from "../../../diary/components/entry-feed/entry-feed";

@Component({
  selector: "app-home",
  imports: [RouterModule, EntryHistory, EntryFeed],
  templateUrl: "./home.html",
  styleUrl: "./home.scss",
})
export class Home {
  private readonly entryApi = inject(EntryApi);
  private readonly diaryEntryApi = inject(DiaryEntryApi);
  private readonly router = inject(Router);
  private readonly limit = 6;

  readonly moodEntriesPage = toSignal(this.entryApi.getPaged(this.limit), {
    initialValue: null,
  });
  readonly diaryEntriesPage = toSignal(
    this.diaryEntryApi.getPaged(this.limit),
    {
      initialValue: null,
    },
  );

  openMoodEntry(entry: Entry): void {
    this.router.navigate(["/mood/entries", entry.id]);
  }

  openDiaryEntry(entry: DiaryEntry): void {
    this.router.navigate(["/diary/entries", entry.id]);
  }
}
