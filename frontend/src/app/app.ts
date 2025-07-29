import { Component, inject, signal } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatRippleModule } from "@angular/material/core";
import { MatToolbarModule } from "@angular/material/toolbar";
import { RouterModule, Router } from "@angular/router";
import { Auth } from "./core/auth/services/auth";
import { toSignal } from "@angular/core/rxjs-interop";

@Component({
  selector: "app-root",
  imports: [RouterModule, MatButtonModule, MatRippleModule, MatToolbarModule],
  templateUrl: "./app.html",
  styleUrl: "./app.scss",
})
export class App {
  protected readonly title = signal("null3");

  private readonly auth = inject(Auth);
  private readonly router = inject(Router);

  readonly user = toSignal(this.auth.user$, { initialValue: null });

  constructor() {
    this.auth.init();
  }

  logout(): void {
    this.auth.logout().subscribe({
      next: () => this.router.navigate(["/login"]),
    });
  }
}
