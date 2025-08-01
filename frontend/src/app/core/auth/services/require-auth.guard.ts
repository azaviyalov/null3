import { Injectable, inject } from "@angular/core";
import { CanActivate, Router, UrlTree } from "@angular/router";
import { Observable } from "rxjs";
import { map, filter } from "rxjs/operators";
import { Auth } from "./auth";

@Injectable({ providedIn: "root" })
export class RequireAuthGuard implements CanActivate {
  private readonly router = inject(Router);
  private readonly auth = inject(Auth);

  canActivate(): Observable<boolean | UrlTree> {
    return this.auth.isAuthenticated$.pipe(
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
