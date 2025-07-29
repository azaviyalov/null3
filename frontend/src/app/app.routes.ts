import { Routes } from "@angular/router";
import { RequireAuthGuard } from "./core/auth/services/require-auth.guard";
import { RequireGuestGuard } from "./core/auth/services/require-guest.guard";

export const routes: Routes = [
  {
    path: "login",
    canActivate: [RequireGuestGuard],
    loadComponent: () =>
      import("./core/auth/pages/login/login").then((m) => m.Login),
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
    redirectTo: "",
  },
];
