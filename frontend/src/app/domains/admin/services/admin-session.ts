import { Injectable, inject } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { BehaviorSubject, Observable } from "rxjs";
import { tap } from "rxjs/operators";
import { environment } from "../../../../environments/environment";
import { AdminLoginRequest } from "../models/admin-login";
import { AdminInviteResponse } from "../models/invite";

@Injectable({ providedIn: "root" })
export class AdminSession {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/admin/auth`;
  private readonly invitesUrl = `${environment.apiUrl}/admin/invites`;

  private readonly _isAuthenticated = new BehaviorSubject<boolean | null>(null);

  init(): void {
    this.http.get<void>(`${this.baseUrl}/me`).subscribe({
      next: () => this._isAuthenticated.next(true),
      error: (error) => {
        console.error(
          "Error during admin authentication initialization:",
          error,
        );
        this._isAuthenticated.next(false);
      },
    });
  }

  get isAuthenticated$(): Observable<boolean | null> {
    return this._isAuthenticated.asObservable();
  }

  login(data: AdminLoginRequest): Observable<void> {
    return this.http
      .post<void>(`${this.baseUrl}/login`, data)
      .pipe(tap(() => this._isAuthenticated.next(true)));
  }

  createInvite(): Observable<AdminInviteResponse> {
    return this.http.post<AdminInviteResponse>(this.invitesUrl, {});
  }

  clearSession(): void {
    this._isAuthenticated.next(false);
  }

  logout(): Observable<void> {
    return this.http
      .post<void>(`${this.baseUrl}/logout`, {})
      .pipe(tap(() => this.clearSession()));
  }
}
