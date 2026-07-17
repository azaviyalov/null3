import { Component, inject } from "@angular/core";
import { RouterModule, Router } from "@angular/router";
import { toSignal } from "@angular/core/rxjs-interop";
import { AdminSession } from "./domains/admin/services/admin-session";
import { UserSession } from "./domains/session/services/user-session";

@Component({
  selector: "app-root",
  imports: [RouterModule],
  templateUrl: "./app.html",
  styleUrl: "./app.scss",
})
export class App {
  private readonly userSession = inject(UserSession);
  private readonly adminSession = inject(AdminSession);
  private readonly router = inject(Router);

  readonly user = toSignal(this.userSession.user$, { initialValue: null });
  readonly adminAuthenticated = toSignal(this.adminSession.isAuthenticated$, {
    initialValue: null,
  });

  constructor() {
    this.userSession.init();
    this.adminSession.init();
  }

  logout(): void {
    this.userSession.logout().subscribe({
      next: () => this.router.navigate(["/login"]),
    });
  }

  logoutAdmin(): void {
    this.adminSession.logout().subscribe({
      next: () => this.router.navigate(["/admin/login"]),
    });
  }
}
