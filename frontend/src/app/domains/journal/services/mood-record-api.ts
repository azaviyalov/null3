import { HttpClient, HttpParams } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import { Observable, map } from "rxjs";
import { environment } from "../../../../environments/environment";
import {
  EditMoodRecordRequest,
  MoodRecordResponse,
  MoodRecord,
} from "../models/mood-record";
import { Page, PageResponse } from "../../../core/utils/page";

@Injectable({
  providedIn: "root",
})
export class MoodRecordApi {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/journal/mood-records`;

  getPaged(
    limit = 10,
    offset = 0,
    deleted = false,
  ): Observable<Page<MoodRecord>> {
    const params = new HttpParams()
      .set("limit", limit)
      .set("offset", offset)
      .set("deleted", deleted);
    return this.http
      .get<PageResponse<MoodRecordResponse>>(this.baseUrl, { params })
      .pipe(map((data) => Page.fromResponse(data, MoodRecord.fromResponse)));
  }

  getById(id: number): Observable<MoodRecord> {
    return this.http
      .get<MoodRecordResponse>(`${this.baseUrl}/${id}`)
      .pipe(map(MoodRecord.fromResponse));
  }

  create(req: EditMoodRecordRequest): Observable<MoodRecord> {
    return this.http
      .post<MoodRecordResponse>(this.baseUrl, req)
      .pipe(map(MoodRecord.fromResponse));
  }

  update(id: number, req: EditMoodRecordRequest): Observable<MoodRecord> {
    return this.http
      .put<MoodRecordResponse>(`${this.baseUrl}/${id}`, req)
      .pipe(map(MoodRecord.fromResponse));
  }

  delete(id: number): Observable<MoodRecord> {
    return this.http
      .delete<MoodRecordResponse>(`${this.baseUrl}/${id}`)
      .pipe(map(MoodRecord.fromResponse));
  }

  restore(id: number): Observable<MoodRecord> {
    return this.http
      .post<MoodRecordResponse>(`${this.baseUrl}/${id}/restore`, {})
      .pipe(map(MoodRecord.fromResponse));
  }
}
