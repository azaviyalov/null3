import { Injectable, inject } from "@angular/core";
import { CanActivate, Router, UrlTree } from "@angular/router";
import { Observable, filter, map } from "rxjs";
import { AdminAuth } from "./admin-auth";

@Injectable({ providedIn: "root" })
export class RequireAdminGuard implements CanActivate {
  private readonly router = inject(Router);
  private readonly adminAuth = inject(AdminAuth);

  canActivate(): Observable<boolean | UrlTree> {
    return this.adminAuth.isAuthenticated$.pipe(
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
