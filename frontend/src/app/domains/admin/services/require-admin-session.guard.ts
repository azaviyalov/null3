import { Injectable, inject } from "@angular/core";
import { CanActivate, Router, UrlTree } from "@angular/router";
import { Observable, filter, map } from "rxjs";
import { AdminSession } from "./admin-session";

@Injectable({ providedIn: "root" })
export class RequireAdminSessionGuard implements CanActivate {
  private readonly router = inject(Router);
  private readonly adminSession = inject(AdminSession);

  canActivate(): Observable<boolean | UrlTree> {
    return this.adminSession.isAuthenticated$.pipe(
      filter((isAuth): isAuth is boolean => isAuth !== null),
      map((isAuth) => {
        if (!isAuth) {
          return this.router.createUrlTree(["/admin/login"]);
        }
        return true;
      }),
    );
  }
}
