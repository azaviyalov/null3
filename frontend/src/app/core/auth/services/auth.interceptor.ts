import { Injectable, inject } from "@angular/core";
import {
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
  HttpErrorResponse,
} from "@angular/common/http";
import { Observable, throwError, ReplaySubject } from "rxjs";
import { catchError, switchMap, take } from "rxjs/operators";
import { Auth } from "./auth";

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  private readonly auth = inject(Auth);
  private refreshInProgressSubject: ReplaySubject<unknown> | null = null;

  intercept(
    req: HttpRequest<unknown>,
    next: HttpHandler,
  ): Observable<HttpEvent<unknown>> {
    // Always send credentials
    const request = req.clone({ withCredentials: true });

    // Only handle API requests
    if (!req.url.includes("/api/")) return next.handle(request);

    return next.handle(request).pipe(
      catchError((err) => this.handleAuthError(err, req, request, next)),
    );
  }

  private handleAuthError(
    err: unknown,
    originalRequest: HttpRequest<unknown>,
    requestWithCredentials: HttpRequest<unknown>,
    next: HttpHandler,
  ): Observable<HttpEvent<unknown>> {
    // Only handle 401 errors for API requests (except refresh endpoint)
    const isAuthError =
      err instanceof HttpErrorResponse &&
      err.status === 401 &&
      !originalRequest.url.endsWith("/auth/refresh") &&
      (originalRequest.url.endsWith("/auth/me") || !!this.auth.currentUser);

    if (!isAuthError) return throwError(() => err);

    // If a refresh is already in progress, reuse it
    if (!this.refreshInProgressSubject) {
      this.refreshInProgressSubject = new ReplaySubject<unknown>(1);
      const subject = this.refreshInProgressSubject;
      this.auth.refresh().pipe(take(1)).subscribe({
        next: (user) => {
          subject.next(user);
          subject.complete();
          this.refreshInProgressSubject = null;
        },
        error: (error) => {
          console.error("Auth refresh failed:", error);
          subject.next(null);
          subject.complete();
          this.refreshInProgressSubject = null;
        }
      });
    }

    return this.refreshInProgressSubject.asObservable().pipe(
      switchMap((user) => {
        if (user) {
          // Retry the original request after successful refresh
          return next.handle(requestWithCredentials);
        }
        // If refresh failed, clear session and propagate error
        this.auth.clearSession();
        return throwError(() => err);
      }),
      catchError(() => {
        this.auth.clearSession();
        return throwError(() => err);
      }),
    );
  }
}
