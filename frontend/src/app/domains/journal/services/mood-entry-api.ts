import { HttpClient, HttpParams } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import { Observable, map } from "rxjs";
import { environment } from "../../../../environments/environment";
import {
  EditMoodEntryRequest,
  MoodEntry,
  MoodEntryResponse,
} from "../models/mood-entry";
import { Page, PageResponse } from "../../../core/utils/page";

@Injectable({
  providedIn: "root",
})
export class MoodEntryApi {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/journal/mood/entries`;

  getPaged(
    limit = 10,
    offset = 0,
    deleted = false,
  ): Observable<Page<MoodEntry>> {
    const params = new HttpParams()
      .set("limit", limit)
      .set("offset", offset)
      .set("deleted", deleted);
    return this.http
      .get<PageResponse<MoodEntryResponse>>(this.baseUrl, { params })
      .pipe(map((data) => Page.fromResponse(data, MoodEntry.fromResponse)));
  }

  getById(id: number): Observable<MoodEntry> {
    return this.http
      .get<MoodEntryResponse>(`${this.baseUrl}/${id}`)
      .pipe(map(MoodEntry.fromResponse));
  }

  create(req: EditMoodEntryRequest): Observable<MoodEntry> {
    return this.http
      .post<MoodEntryResponse>(this.baseUrl, req)
      .pipe(map(MoodEntry.fromResponse));
  }

  update(id: number, req: EditMoodEntryRequest): Observable<MoodEntry> {
    return this.http
      .put<MoodEntryResponse>(`${this.baseUrl}/${id}`, req)
      .pipe(map(MoodEntry.fromResponse));
  }

  delete(id: number): Observable<MoodEntry> {
    return this.http
      .delete<MoodEntryResponse>(`${this.baseUrl}/${id}`)
      .pipe(map(MoodEntry.fromResponse));
  }

  restore(id: number): Observable<MoodEntry> {
    return this.http
      .post<MoodEntryResponse>(`${this.baseUrl}/${id}/restore`, {})
      .pipe(map(MoodEntry.fromResponse));
  }
}
