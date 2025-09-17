import { Injectable, inject, signal } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { Observable, tap } from "rxjs";
import { environment } from "../../../../environments/environment";

export interface AdminLoginRequest {
  username: string;
  password: string;
}

export interface User {
  id: number;
  email: string;
  name: string;
}

export interface CreateUserRequest {
  email: string;
  name: string;
  password: string;
}

export interface UpdateUserRequest {
  email?: string;
  name?: string;
  password?: string;
}

export interface RefreshToken {
  id: number;
  user_id: number;
  value: string;
  created_at: string;
  expires_at: string;
}

@Injectable({
  providedIn: "root",
})
export class AdminAuth {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = environment.apiUrl || "http://localhost:8080/api";

  readonly isLoggedIn = signal(false);

  login(req: AdminLoginRequest): Observable<void> {
    return this.http.post<void>(`${this.baseUrl}/admin/login`, req).pipe(
      tap(() => {
        this.isLoggedIn.set(true);
      })
    );
  }

  logout(): Observable<void> {
    return this.http.post<void>(`${this.baseUrl}/admin/logout`, {}).pipe(
      tap(() => {
        this.isLoggedIn.set(false);
      })
    );
  }

  // User management methods
  getUsers(): Observable<User[]> {
    return this.http.get<User[]>(`${this.baseUrl}/admin/users`);
  }

  createUser(req: CreateUserRequest): Observable<User> {
    return this.http.post<User>(`${this.baseUrl}/admin/users`, req);
  }

  getUser(id: number): Observable<User> {
    return this.http.get<User>(`${this.baseUrl}/admin/users/${id}`);
  }

  updateUser(id: number, req: UpdateUserRequest): Observable<User> {
    return this.http.put<User>(`${this.baseUrl}/admin/users/${id}`, req);
  }

  deleteUser(id: number): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/admin/users/${id}`);
  }

  // Refresh token management
  getRefreshTokens(): Observable<RefreshToken[]> {
    return this.http.get<RefreshToken[]>(`${this.baseUrl}/admin/refresh-tokens`);
  }

  deleteRefreshToken(value: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/admin/refresh-tokens/${value}`);
  }
}