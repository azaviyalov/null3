import { HttpClient, HttpParams } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import { Observable, map } from "rxjs";
import { environment } from "../../../../environments/environment";
import {
  PaginatedEntryList,
  EntryView,
  Entry,
  PaginatedEntryListView,
  EditEntryRequest,
} from "../models/entry";

@Injectable({
  providedIn: "root",
})
export class EntryApi {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/mood/entries`;

  getPaged(
    limit = 10,
    offset = 0,
    deleted = false,
  ): Observable<PaginatedEntryList> {
    const params = new HttpParams()
      .set("limit", limit)
      .set("offset", offset)
      .set("deleted", deleted);

    return this.http
      .get<PaginatedEntryListView>(this.baseUrl, { params })
      .pipe(map(PaginatedEntryList.fromView));
  }

  getById(id: number): Observable<Entry> {
    return this.http
      .get<EntryView>(`${this.baseUrl}/${id}`)
      .pipe(map(Entry.fromView));
  }

  create(req: EditEntryRequest): Observable<Entry> {
    return this.http
      .post<EntryView>(this.baseUrl, req)
      .pipe(map(Entry.fromView));
  }

  update(id: number, req: EditEntryRequest): Observable<Entry> {
    return this.http
      .put<EntryView>(`${this.baseUrl}/${id}`, req)
      .pipe(map(Entry.fromView));
  }

  delete(id: number): Observable<Entry> {
    return this.http
      .delete<EntryView>(`${this.baseUrl}/${id}`)
      .pipe(map(Entry.fromView));
  }

  restore(id: number): Observable<Entry> {
    return this.http
      .post<EntryView>(`${this.baseUrl}/${id}/restore`, {})
      .pipe(map(Entry.fromView));
  }
}
