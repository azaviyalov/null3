import { Component, inject } from "@angular/core";
import { Router, RouterModule } from "@angular/router";
import { toSignal } from "@angular/core/rxjs-interop";
import { MoodRecordHistory } from "../../../journal/components/mood-record-history";
import { MoodRecordApi } from "../../../journal/services/mood-record-api";
import { DiaryEntryApi } from "../../../journal/services/diary-entry-api";
import { DiaryEntryFeed } from "../../../journal/components/diary-entry-feed";
import { DiaryEntry } from "../../../journal/models/diary-entry";
import { MoodRecord } from "../../../journal/models/mood-record";

@Component({
  selector: "app-home",
  imports: [RouterModule, MoodRecordHistory, DiaryEntryFeed],
  templateUrl: "./home.html",
  styleUrl: "./home.scss",
})
export class Home {
  private readonly moodRecordApi = inject(MoodRecordApi);
  private readonly diaryEntryApi = inject(DiaryEntryApi);
  private readonly router = inject(Router);
  private readonly limit = 6;

  readonly moodRecordsPage = toSignal(this.moodRecordApi.getPaged(this.limit), {
    initialValue: null,
  });
  readonly diaryEntriesPage = toSignal(
    this.diaryEntryApi.getPaged(this.limit),
    {
      initialValue: null,
    },
  );

  openMoodRecord(entry: MoodRecord): void {
    this.router.navigate(["/mood/records", entry.id]);
  }

  openDiaryEntry(entry: DiaryEntry): void {
    this.router.navigate(["/diary/entries", entry.id]);
  }
}
