import { Injectable, inject } from "@angular/core";
import { CanActivate, Router, UrlTree } from "@angular/router";
import { UserSession } from "./user-session";
import { Observable, map, filter } from "rxjs";

@Injectable({ providedIn: "root" })
export class RequireUserGuestGuard implements CanActivate {
  private readonly userSession = inject(UserSession);
  private readonly router = inject(Router);

  canActivate(): Observable<boolean | UrlTree> {
    return this.userSession.isAuthenticated$.pipe(
      filter((isAuth): isAuth is boolean => isAuth !== null),
      map((isAuth) => {
        if (isAuth) {
          return this.router.createUrlTree(["/"]);
        }
        return true;
      }),
    );
  }
}
