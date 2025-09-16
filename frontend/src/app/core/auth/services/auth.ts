import { Injectable, inject } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { BehaviorSubject, Observable, of } from "rxjs";
import { tap, catchError } from "rxjs/operators";
import { environment } from "../../../../environments/environment";
import { LoginRequest } from "../models/login";
import { UserResponse } from "../models/user";

@Injectable({ providedIn: "root" })
export class Auth {
  private readonly http = inject(HttpClient);
  private readonly loginUrl = `${environment.apiUrl}/auth/login`;
  private readonly meUrl = `${environment.apiUrl}/auth/me`;

  private _user = new BehaviorSubject<UserResponse | null>(null);
  private _isAuthenticated = new BehaviorSubject<boolean | null>(null);

  init(): void {
    this.http
      .get<UserResponse>(this.meUrl)
      .pipe(
        catchError((error) => {
          console.error("Error during authentication initialization:", error);
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

  login(data: LoginRequest): Observable<UserResponse> {
    return this.http.post<UserResponse>(this.loginUrl, data).pipe(
      tap((response) => {
        this._user.next(response);
        this._isAuthenticated.next(true);
      }),
    );
  }

  clearSession() {
    this._user.next(null);
    this._isAuthenticated.next(false);
  }

  logout(): Observable<void> {
    return this.http.post<void>(`${environment.apiUrl}/auth/logout`, {}).pipe(
      tap(() => {
        this._user.next(null);
        this._isAuthenticated.next(false);
      }),
    );
  }

  refresh(): Observable<UserResponse | null> {
    return this.http
      .post<UserResponse>(`${environment.apiUrl}/auth/refresh`, {})
      .pipe(
        tap((response) => {
          this._user.next(response);
          this._isAuthenticated.next(true);
        }),
        catchError(() => {
          this.clearSession();
          return of(null);
        }),
      );
  }
}
