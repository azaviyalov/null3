import { Component, inject, signal } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatRippleModule } from "@angular/material/core";
import { MatToolbarModule } from "@angular/material/toolbar";
import { RouterModule, Router } from "@angular/router";
import { Auth } from "./core/auth/services/auth";
import { toSignal } from "@angular/core/rxjs-interop";
import { AdminAuth } from "./core/admin/services/admin-auth";

@Component({
  selector: "app-root",
  imports: [RouterModule, MatButtonModule, MatRippleModule, MatToolbarModule],
  templateUrl: "./app.html",
  styleUrl: "./app.scss",
})
export class App {
  protected readonly title = signal("null3");

  private readonly auth = inject(Auth);
  private readonly adminAuth = inject(AdminAuth);
  private readonly router = inject(Router);

  readonly user = toSignal(this.auth.user$, { initialValue: null });
  readonly adminUser = toSignal(this.adminAuth.user$, { initialValue: null });

  constructor() {
    this.auth.init();
    this.adminAuth.init();
  }

  logout(): void {
    this.auth.logout().subscribe({
      next: () => this.router.navigate(["/login"]),
    });
  }

  logoutAdmin(): void {
    this.adminAuth.logout().subscribe({
      next: () => this.router.navigate(["/admin/login"]),
    });
  }
}
