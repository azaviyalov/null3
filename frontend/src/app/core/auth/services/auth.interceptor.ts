import { Injectable, inject } from "@angular/core";
import {
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
  HttpErrorResponse,
} from "@angular/common/http";
import { Observable, throwError } from "rxjs";
import { catchError } from "rxjs/operators";
import { Auth } from "./auth";

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  private readonly auth = inject(Auth);

  intercept(
    req: HttpRequest<unknown>,
    next: HttpHandler,
  ): Observable<HttpEvent<unknown>> {
    return next
      .handle(
        req.clone({
          withCredentials: true,
        }),
      )
      .pipe(
        catchError((err) => {
          if (err instanceof HttpErrorResponse && err.status === 401) {
            this.auth.clearSession();
          }
          return throwError(() => err);
        }),
      );
  }
}
