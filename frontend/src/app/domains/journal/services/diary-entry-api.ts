import { HttpClient, HttpParams } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import { Observable, map } from "rxjs";
import { environment } from "../../../../environments/environment";
import {
  DiaryEntry,
  DiaryEntryResponse,
  EditDiaryEntryRequest,
} from "../models/diary-entry";
import { Page, PageResponse } from "../../../core/utils/page";

@Injectable({
  providedIn: "root",
})
export class DiaryEntryApi {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/journal/diary/entries`;

  getPaged(
    limit = 10,
    offset = 0,
    deleted = false,
  ): Observable<Page<DiaryEntry>> {
    const params = new HttpParams()
      .set("limit", limit)
      .set("offset", offset)
      .set("deleted", deleted);

    return this.http
      .get<PageResponse<DiaryEntryResponse>>(this.baseUrl, { params })
      .pipe(map((data) => Page.fromResponse(data, DiaryEntry.fromResponse)));
  }

  getById(id: number): Observable<DiaryEntry> {
    return this.http
      .get<DiaryEntryResponse>(`${this.baseUrl}/${id}`)
      .pipe(map(DiaryEntry.fromResponse));
  }

  create(req: EditDiaryEntryRequest): Observable<DiaryEntry> {
    return this.http
      .post<DiaryEntryResponse>(this.baseUrl, req)
      .pipe(map(DiaryEntry.fromResponse));
  }

  update(id: number, req: EditDiaryEntryRequest): Observable<DiaryEntry> {
    return this.http
      .put<DiaryEntryResponse>(`${this.baseUrl}/${id}`, req)
      .pipe(map(DiaryEntry.fromResponse));
  }

  delete(id: number): Observable<DiaryEntry> {
    return this.http
      .delete<DiaryEntryResponse>(`${this.baseUrl}/${id}`)
      .pipe(map(DiaryEntry.fromResponse));
  }

  restore(id: number): Observable<DiaryEntry> {
    return this.http
      .post<DiaryEntryResponse>(`${this.baseUrl}/${id}/restore`, {})
      .pipe(map(DiaryEntry.fromResponse));
  }
}
