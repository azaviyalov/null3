import { Routes } from "@angular/router";
import { RequireAdminSessionGuard } from "./domains/admin/services/require-admin-session.guard";
import { RequireAdminGuestGuard } from "./domains/admin/services/require-admin-guest.guard";
import { RequireUserSessionGuard } from "./domains/session/services/require-user-session.guard";
import { RequireUserGuestGuard } from "./domains/session/services/require-user-guest.guard";

export const routes: Routes = [
  {
    path: "login",
    canActivate: [RequireUserGuestGuard],
    loadComponent: () =>
      import("./domains/session/pages/login/login").then((m) => m.Login),
  },
  {
    path: "forgot-password",
    canActivate: [RequireUserGuestGuard],
    loadComponent: () =>
      import("./domains/account/pages/forgot-password/forgot-password").then(
        (m) => m.ForgotPassword,
      ),
  },
  {
    path: "reset-password",
    canActivate: [RequireUserGuestGuard],
    loadComponent: () =>
      import("./domains/account/pages/reset-password/reset-password").then(
        (m) => m.ResetPassword,
      ),
  },
  {
    path: "invite/:token",
    canActivate: [RequireUserGuestGuard],
    loadComponent: () =>
      import("./domains/account/pages/invite-register/invite-register").then(
        (m) => m.InviteRegister,
      ),
  },
  {
    path: "admin/login",
    canActivate: [RequireAdminGuestGuard],
    loadComponent: () =>
      import("./domains/admin/pages/admin-login/admin-login").then(
        (m) => m.AdminLogin,
      ),
  },
  {
    path: "admin/invites",
    canActivate: [RequireAdminSessionGuard],
    loadComponent: () =>
      import("./domains/admin/pages/admin-invites/admin-invites").then(
        (m) => m.AdminInvites,
      ),
  },
  { path: "logout", redirectTo: "/login" },
  {
    path: "",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/dashboard/pages/home/home").then((m) => m.Home),
  },
  {
    path: "mood/entries",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/mood/pages/entry-list/entry-list").then(
        (m) => m.EntryList,
      ),
  },
  {
    path: "mood/entries/create",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/mood/pages/entry-create/entry-create").then(
        (m) => m.EntryCreate,
      ),
  },
  {
    path: "mood/entries/:id",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/mood/pages/entry-view/entry-view").then(
        (m) => m.EntryView,
      ),
  },
  {
    path: "mood/entries/:id/update",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/mood/pages/entry-update/entry-update").then(
        (m) => m.EntryUpdate,
      ),
  },
  {
    path: "diary/entries",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/diary/pages/entry-list/entry-list").then(
        (m) => m.EntryList,
      ),
  },
  {
    path: "diary/entries/create",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/diary/pages/entry-create/entry-create").then(
        (m) => m.EntryCreate,
      ),
  },
  {
    path: "diary/entries/:id",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/diary/pages/entry-view/entry-view").then(
        (m) => m.EntryView,
      ),
  },
  {
    path: "diary/entries/:id/update",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/diary/pages/entry-update/entry-update").then(
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
