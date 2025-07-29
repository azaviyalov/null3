import { Injectable, inject } from "@angular/core";
import { CanActivate, Router, UrlTree } from "@angular/router";
import { Auth } from "./auth";
import { Observable, map } from "rxjs";

@Injectable({ providedIn: "root" })
export class RequireGuestGuard implements CanActivate {
  private readonly auth = inject(Auth);
  private readonly router = inject(Router);

  canActivate(): Observable<boolean | UrlTree> {
    return this.auth.isAuthenticated$.pipe(
      map((isAuth) => {
        if (isAuth) {
          return this.router.createUrlTree(["/"]);
        }
        return true;
      }),
    );
  }
}
