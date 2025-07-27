import { Routes } from "@angular/router";

export const routes: Routes = [
  {
    path: "",
    loadComponent: () => import("./pages/home/home").then((m) => m.Home),
  },
  {
    path: "mood/entries",
    loadComponent: () =>
      import("./domains/mood/pages/entry-list/entry-list").then(
        (m) => m.EntryList,
      ),
  },
  {
    path: "mood/entries/create",
    loadComponent: () =>
      import("./domains/mood/pages/entry-create/entry-create").then(
        (m) => m.EntryCreate,
      ),
  },
  {
    path: "mood/entries/:id",
    loadComponent: () =>
      import("./domains/mood/pages/entry-view/entry-view").then(
        (m) => m.EntryView,
      ),
  },
  {
    path: "mood/entries/:id/update",
    loadComponent: () =>
      import("./domains/mood/pages/entry-update/entry-update").then(
        (m) => m.EntryUpdate,
      ),
  },
  {
    path: "about",
    loadComponent: () => import("./pages/about/about").then((m) => m.About),
  },
  {
    path: "**",
    redirectTo: "",
  },
];
