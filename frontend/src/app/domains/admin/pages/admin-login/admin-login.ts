import { Component, inject, signal } from "@angular/core";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { HttpErrorResponse } from "@angular/common/http";
import { Router, RouterModule } from "@angular/router";
import { AdminSession } from "../../services/admin-session";
import { AdminLoginRequest } from "../../models/admin-login";

const ADMIN_HOME_ROUTE = "/admin/invites";

@Component({
  selector: "app-admin-login",
  standalone: true,
  imports: [ReactiveFormsModule, RouterModule],
  templateUrl: "./admin-login.html",
  styleUrl: "./admin-login.scss",
})
export class AdminLogin {
  private readonly adminSession = inject(AdminSession);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  readonly form = this.fb.group({
    password: ["", Validators.required],
  });

  readonly error = signal<string | null>(null);
  readonly isSubmitting = signal(false);

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.error.set(null);
    this.isSubmitting.set(true);

    const req: AdminLoginRequest = {
      password: this.form.value.password!,
    };

    this.adminSession.login(req).subscribe({
      next: () => {
        this.isSubmitting.set(false);
        void this.router.navigate([ADMIN_HOME_ROUTE]);
      },
      error: (error) => this.handleError(error),
    });
  }

  private handleError(error: HttpErrorResponse): void {
    this.isSubmitting.set(false);
    if (error.status === 0) {
      this.error.set("Network error. Please check your connection.");
      return;
    }
    if (error.status === 401) {
      this.error.set("Incorrect admin credentials.");
      return;
    }
    this.error.set("Failed to sign in to the admin area.");
  }
}
