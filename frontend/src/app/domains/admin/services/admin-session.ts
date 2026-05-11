import { Injectable, inject } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { BehaviorSubject, Observable, of } from "rxjs";
import { catchError, tap } from "rxjs/operators";
import { environment } from "../../../../environments/environment";
import { LoginRequest } from "../../session/models/login";
import { UserResponse } from "../../session/models/user";
import { AdminInviteResponse } from "../models/invite";

@Injectable({ providedIn: "root" })
export class AdminSession {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/admin/auth`;
  private readonly invitesUrl = `${environment.apiUrl}/admin/invites`;

  private readonly _user = new BehaviorSubject<UserResponse | null>(null);
  private readonly _isAuthenticated = new BehaviorSubject<boolean | null>(null);

  init(): void {
    this.http
      .get<UserResponse>(`${this.baseUrl}/me`)
      .pipe(
        catchError((error) => {
          console.error(
            "Error during admin authentication initialization:",
            error,
          );
          return of(null);
        }),
      )
      .subscribe({
        next: (user) => {
          this._user.next(user);
          this._isAuthenticated.next(user !== null);
        },
      });
  }

  get user$(): Observable<UserResponse | null> {
    return this._user.asObservable();
  }

  get isAuthenticated$(): Observable<boolean | null> {
    return this._isAuthenticated.asObservable();
  }

  get currentUser(): UserResponse | null {
    return this._user.getValue();
  }

  setAuthenticatedUser(user: UserResponse): void {
    this._user.next(user);
    this._isAuthenticated.next(true);
  }

  login(data: LoginRequest): Observable<UserResponse> {
    return this.http
      .post<UserResponse>(`${this.baseUrl}/login`, data)
      .pipe(tap((response) => this.setAuthenticatedUser(response)));
  }

  createInvite(): Observable<AdminInviteResponse> {
    return this.http.post<AdminInviteResponse>(this.invitesUrl, {});
  }

  clearSession(): void {
    this._user.next(null);
    this._isAuthenticated.next(false);
  }

  logout(): Observable<void> {
    return this.http
      .post<void>(`${this.baseUrl}/logout`, {})
      .pipe(tap(() => this.clearSession()));
  }

  refresh(): Observable<UserResponse | null> {
    return this.http.post<UserResponse>(`${this.baseUrl}/refresh`, {}).pipe(
      tap((response) => this.setAuthenticatedUser(response)),
      catchError(() => {
        this.clearSession();
        return of(null);
      }),
    );
  }
}
