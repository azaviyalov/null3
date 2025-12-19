import { HttpClient, HttpParams } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import { Observable, map } from "rxjs";
import { environment } from "../../../../environments/environment";
import {
  EntryResponse,
  Entry,
  EditEntryRequest,
} from "../models/entry";
import { Page, PageResponse } from "../../../core/utils/page";
import { Auth } from "../../../core/auth/services/auth";

@Injectable({
  providedIn: "root",
})
export class EntryApi {
  private readonly http = inject(HttpClient);
  private readonly auth = inject(Auth);
  private readonly baseUrl = `${environment.apiUrl}/mood/entries`;

  getPaged(
    limit = 10,
    offset = 0,
    deleted = false,
  ): Observable<Page<Entry>> {
    const params = new HttpParams()
      .set("limit", limit)
      .set("offset", offset)
      .set("deleted", deleted);
    return this.http
      .get<PageResponse<EntryResponse>>(this.baseUrl, { params })
      .pipe(map(data => Page.fromResponse(data, Entry.fromResponse)));
  }

  getById(id: number): Observable<Entry> {
    return this.http
      .get<EntryResponse>(`${this.baseUrl}/${id}`)
      .pipe(map(Entry.fromResponse));
  }

  create(req: EditEntryRequest): Observable<Entry> {
    return this.http
      .post<EntryResponse>(this.baseUrl, req)
      .pipe(map(Entry.fromResponse));
  }

  update(id: number, req: EditEntryRequest): Observable<Entry> {
    return this.http
      .put<EntryResponse>(`${this.baseUrl}/${id}`, req)
      .pipe(map(Entry.fromResponse));
  }

  delete(id: number): Observable<Entry> {
    return this.http
      .delete<EntryResponse>(`${this.baseUrl}/${id}`)
      .pipe(map(Entry.fromResponse));
  }

  restore(id: number): Observable<Entry> {
    return this.http
      .post<EntryResponse>(`${this.baseUrl}/${id}/restore`, {})
      .pipe(map(Entry.fromResponse));
  }
}
