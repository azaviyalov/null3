import { Injectable, inject } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { BehaviorSubject, Observable, of } from "rxjs";
import { tap, catchError, filter, ignoreElements } from "rxjs/operators";
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
    this.fetchCurrentUser().subscribe({
      next: (user) => {
        this._user.next(user);
        this._isAuthenticated.next(user !== null);
      },
      error: () => {
        this._user.next(null);
        this._isAuthenticated.next(false);
      },
    });
  }

  get user$(): Observable<UserResponse | null> {
    return this._user.asObservable();
  }

  get isAuthenticated$(): Observable<boolean> {
    return this._isAuthenticated
      .asObservable()
      .pipe(filter((value) => value !== null));
  }

  login(data: LoginRequest): Observable<void> {
    return this.http.post<UserResponse>(this.loginUrl, data).pipe(
      tap((response) => {
        this._user.next(response);
        this._isAuthenticated.next(true);
      }),
      ignoreElements(),
    );
  }

  logout(): void {
    this._user.next(null);
    this._isAuthenticated.next(false);
  }
  private fetchCurrentUser(): Observable<UserResponse | null> {
    return this.http
      .get<UserResponse>(this.meUrl)
      .pipe(catchError(() => of(null)));
  }
}
