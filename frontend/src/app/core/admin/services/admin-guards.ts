import { inject } from "@angular/core";
import { CanActivateFn, Router } from "@angular/router";
import { AdminAuth } from "../services/admin-auth";

export const RequireAdminGuard: CanActivateFn = () => {
  const adminAuth = inject(AdminAuth);
  const router = inject(Router);

  if (adminAuth.isLoggedIn()) {
    return true;
  }

  router.navigate(["/admin/login"]);
  return false;
};

export const RequireAdminGuestGuard: CanActivateFn = () => {
  const adminAuth = inject(AdminAuth);
  const router = inject(Router);

  if (!adminAuth.isLoggedIn()) {
    return true;
  }

  router.navigate(["/admin"]);
  return false;
};