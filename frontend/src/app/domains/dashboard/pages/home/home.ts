import { Component, inject } from "@angular/core";
import { Router, RouterModule } from "@angular/router";
import { toSignal } from "@angular/core/rxjs-interop";
import { MoodEntryHistory } from "../../../journal/components/mood-entry-history";
import { MoodEntryApi } from "../../../journal/services/mood-entry-api";
import { DiaryEntryApi } from "../../../journal/services/diary-entry-api";
import { DiaryEntryFeed } from "../../../journal/components/diary-entry-feed";
import { DiaryEntry } from "../../../journal/models/diary-entry";
import { MoodEntry } from "../../../journal/models/mood-entry";

@Component({
  selector: "app-home",
  imports: [RouterModule, MoodEntryHistory, DiaryEntryFeed],
  templateUrl: "./home.html",
  styleUrl: "./home.scss",
})
export class Home {
  private readonly moodEntryApi = inject(MoodEntryApi);
  private readonly diaryEntryApi = inject(DiaryEntryApi);
  private readonly router = inject(Router);
  private readonly limit = 6;

  readonly moodEntriesPage = toSignal(this.moodEntryApi.getPaged(this.limit), {
    initialValue: null,
  });
  readonly diaryEntriesPage = toSignal(
    this.diaryEntryApi.getPaged(this.limit),
    {
      initialValue: null,
    },
  );

  openMoodEntry(entry: MoodEntry): void {
    this.router.navigate(["/mood/entries", entry.id]);
  }

  openDiaryEntry(entry: DiaryEntry): void {
    this.router.navigate(["/diary/entries", entry.id]);
  }
}
