import { Injectable, inject } from "@angular/core";
import { CanActivate, Router, UrlTree } from "@angular/router";
import { Observable } from "rxjs";
import { map, filter } from "rxjs/operators";
import { UserSession } from "./user-session";

@Injectable({ providedIn: "root" })
export class RequireUserSessionGuard implements CanActivate {
  private readonly router = inject(Router);
  private readonly userSession = inject(UserSession);

  canActivate(): Observable<boolean | UrlTree> {
    return this.userSession.isAuthenticated$.pipe(
      filter((isAuth): isAuth is boolean => isAuth !== null),
      map((isAuth) => {
        if (!isAuth) {
          return this.router.createUrlTree(["/login"]);
        }
        return true;
      }),
    );
  }
}
