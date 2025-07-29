import { Injectable, inject } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { BehaviorSubject, Observable, of } from "rxjs";
import { map, tap, catchError, filter } from "rxjs/operators";
import { environment } from "../../../../environments/environment";
import { LoginRequest } from "../models/login";
import { UserInfo } from "../models/user-info";

@Injectable({ providedIn: "root" })
export class Auth {
  private readonly http = inject(HttpClient);
  private readonly loginUrl = `${environment.apiUrl}/auth/login`;
  private readonly meUrl = `${environment.apiUrl}/auth/me`;

  private _user = new BehaviorSubject<UserInfo | null>(null);
  private _isAuthenticated = new BehaviorSubject<boolean | null>(null);

  init(): void {
    this.fetchUser().subscribe({
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

  get user$(): Observable<UserInfo | null> {
    return this._user.asObservable();
  }

  get isAuthenticated$(): Observable<boolean> {
    return this._isAuthenticated
      .asObservable()
      .pipe(filter((value) => value !== null));
  }

  login(data: LoginRequest): Observable<void> {
    return this.http.post<UserInfo>(this.loginUrl, data).pipe(
      tap((response) => {
        localStorage.setItem("jwt_token", response.token);
        this._user.next(response);
        this._isAuthenticated.next(true);
      }),
      map(() => void 0),
    );
  }

  logout(): void {
    localStorage.removeItem("jwt_token");
    this._user.next(null);
    this._isAuthenticated.next(false);
  }

  getToken(): string | null {
    return localStorage.getItem("jwt_token");
  }

  private fetchUser(): Observable<UserInfo | null> {
    const token = this.getToken();
    if (!token) {
      return of(null);
    }
    return this.http.get<UserInfo>(this.meUrl).pipe(catchError(() => of(null)));
  }
}
