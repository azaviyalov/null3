import { Injectable, inject } from "@angular/core";
import {
  HttpErrorResponse,
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from "@angular/common/http";
import { Router } from "@angular/router";
import { Observable, throwError } from "rxjs";
import { catchError } from "rxjs/operators";
import { AdminSession } from "./admin-session";

@Injectable()
export class AdminSessionInterceptor implements HttpInterceptor {
  private readonly adminSession = inject(AdminSession);
  private readonly router = inject(Router);

  intercept(
    req: HttpRequest<unknown>,
    next: HttpHandler,
  ): Observable<HttpEvent<unknown>> {
    if (!req.url.includes("/api/admin/")) return next.handle(req);

    const request = req.clone({ withCredentials: true });

    return next.handle(request).pipe(
      catchError((error) => {
        if (error instanceof HttpErrorResponse && error.status === 401) {
          this.adminSession.clearSession();
          void this.router.navigate(["/admin/login"]);
        }
        return throwError(() => error);
      }),
    );
  }
}
