import { Injectable, inject } from "@angular/core";
import {
  HttpErrorResponse,
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from "@angular/common/http";
import { Observable, ReplaySubject, throwError } from "rxjs";
import { catchError, switchMap, take } from "rxjs/operators";
import { AdminSession } from "./admin-session";

@Injectable()
export class AdminSessionInterceptor implements HttpInterceptor {
  private readonly adminSession = inject(AdminSession);

  private refreshInProgress: ReplaySubject<unknown> | null = null;

  intercept(
    req: HttpRequest<unknown>,
    next: HttpHandler,
  ): Observable<HttpEvent<unknown>> {
    if (!this.isAdminApiRequest(req.url)) return next.handle(req);

    const request = req.clone({ withCredentials: true });

    return next
      .handle(request)
      .pipe(catchError((err) => this.handleAuthError(err, req, request, next)));
  }

  private handleAuthError(
    err: unknown,
    originalRequest: HttpRequest<unknown>,
    requestWithCredentials: HttpRequest<unknown>,
    next: HttpHandler,
  ): Observable<HttpEvent<unknown>> {
    const isHttpError = err instanceof HttpErrorResponse;
    const isUnauthorized = isHttpError && err.status === 401;
    const isRefreshEndpoint = originalRequest.url.endsWith(
      "/admin/auth/refresh",
    );
    const isMeEndpoint = originalRequest.url.endsWith("/admin/auth/me");
    const hasCurrentUser = !!this.adminSession.currentUser;

    const isAuthError =
      isUnauthorized && !isRefreshEndpoint && (isMeEndpoint || hasCurrentUser);

    if (!isAuthError) return throwError(() => err);

    const subject = this.getOrStartRefresh();

    return subject.asObservable().pipe(
      take(1),
      switchMap((user) => {
        if (user) {
          return next.handle(requestWithCredentials);
        }

        this.adminSession.clearSession();
        return throwError(() => err);
      }),
      catchError(() => {
        this.adminSession.clearSession();
        return throwError(() => err);
      }),
    );
  }

  private isAdminApiRequest(url: string): boolean {
    return url.includes("/api/admin/");
  }

  private getOrStartRefresh(): ReplaySubject<unknown> {
    if (this.refreshInProgress) return this.refreshInProgress;

    const subject = new ReplaySubject<unknown>(1);
    this.refreshInProgress = subject;

    this.adminSession
      .refresh()
      .pipe(take(1))
      .subscribe({
        next: (user) => this.finishRefresh(subject, user),
        error: () => this.finishRefresh(subject, null),
      });

    return subject;
  }

  private finishRefresh(subject: ReplaySubject<unknown>, user: unknown): void {
    subject.next(user);
    subject.complete();
    this.refreshInProgress = null;
  }
}
