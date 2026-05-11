import { Routes } from "@angular/router";
import { RequireAuthGuard } from "./core/auth/services/require-auth.guard";
import { RequireGuestGuard } from "./core/auth/services/require-guest.guard";
import { RequireAdminGuard } from "./core/admin/services/require-admin.guard";
import { RequireAdminGuestGuard } from "./core/admin/services/require-admin-guest.guard";

export const routes: Routes = [
  {
    path: "login",
    canActivate: [RequireGuestGuard],
    loadComponent: () =>
      import("./core/auth/pages/login/login").then((m) => m.Login),
  },
  {
    path: "forgot-password",
    canActivate: [RequireGuestGuard],
    loadComponent: () =>
      import("./core/auth/pages/forgot-password/forgot-password").then(
        (m) => m.ForgotPassword,
      ),
  },
  {
    path: "reset-password",
    canActivate: [RequireGuestGuard],
    loadComponent: () =>
      import("./core/auth/pages/reset-password/reset-password").then(
        (m) => m.ResetPassword,
      ),
  },
  {
    path: "invite/:token",
    canActivate: [RequireGuestGuard],
    loadComponent: () =>
      import("./core/auth/pages/invite-register/invite-register").then(
        (m) => m.InviteRegister,
      ),
  },
  {
    path: "admin/login",
    canActivate: [RequireAdminGuestGuard],
    loadComponent: () =>
      import("./core/admin/pages/admin-login/admin-login").then(
        (m) => m.AdminLogin,
      ),
  },
  {
    path: "admin/invites",
    canActivate: [RequireAdminGuard],
    loadComponent: () =>
      import("./core/admin/pages/admin-invites/admin-invites").then(
        (m) => m.AdminInvites,
      ),
  },
  { path: "logout", redirectTo: "/login" },
  {
    path: "",
    canActivate: [RequireAuthGuard],
    loadComponent: () => import("./core/pages/home/home").then((m) => m.Home),
  },
  {
    path: "mood/entries",
    canActivate: [RequireAuthGuard],
    loadComponent: () =>
      import("./domains/mood/pages/entry-list/entry-list").then(
        (m) => m.EntryList,
      ),
  },
  {
    path: "mood/entries/create",
    canActivate: [RequireAuthGuard],
    loadComponent: () =>
      import("./domains/mood/pages/entry-create/entry-create").then(
        (m) => m.EntryCreate,
      ),
  },
  {
    path: "mood/entries/:id",
    canActivate: [RequireAuthGuard],
    loadComponent: () =>
      import("./domains/mood/pages/entry-view/entry-view").then(
        (m) => m.EntryView,
      ),
  },
  {
    path: "mood/entries/:id/update",
    canActivate: [RequireAuthGuard],
    loadComponent: () =>
      import("./domains/mood/pages/entry-update/entry-update").then(
        (m) => m.EntryUpdate,
      ),
  },
  {
    path: "about",
    loadComponent: () =>
      import("./core/pages/about/about").then((m) => m.About),
  },
  {
    path: "**",
    redirectTo: "/",
  },
];
