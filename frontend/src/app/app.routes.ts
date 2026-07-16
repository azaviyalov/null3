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
    path: "mood-records",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/journal/pages/mood-record-list").then(
        (m) => m.MoodRecordList,
      ),
  },
  {
    path: "mood-records/create",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/journal/pages/mood-record-create").then(
        (m) => m.MoodRecordCreate,
      ),
  },
  {
    path: "mood-records/:id",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/journal/pages/mood-record-view").then(
        (m) => m.MoodRecordView,
      ),
  },
  {
    path: "mood-records/:id/update",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/journal/pages/mood-record-update").then(
        (m) => m.MoodRecordUpdate,
      ),
  },
  {
    path: "diary-entries",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/journal/pages/diary-entry-list").then(
        (m) => m.DiaryEntryList,
      ),
  },
  {
    path: "diary-entries/create",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/journal/pages/diary-entry-create").then(
        (m) => m.DiaryEntryCreate,
      ),
  },
  {
    path: "diary-entries/:id",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/journal/pages/diary-entry-view").then(
        (m) => m.DiaryEntryView,
      ),
  },
  {
    path: "diary-entries/:id/update",
    canActivate: [RequireUserSessionGuard],
    loadComponent: () =>
      import("./domains/journal/pages/diary-entry-update").then(
        (m) => m.DiaryEntryUpdate,
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
