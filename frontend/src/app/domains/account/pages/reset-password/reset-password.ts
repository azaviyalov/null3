import { Component, inject, signal } from "@angular/core";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { HttpErrorResponse } from "@angular/common/http";
import { MatButtonModule } from "@angular/material/button";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { ActivatedRoute, Router, RouterModule } from "@angular/router";
import { AccountApi } from "../../services/account-api";
import { ResetPasswordRequest } from "../../models/password";

const LOGIN_ROUTE = "/login";

@Component({
  selector: "app-reset-password",
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    RouterModule,
  ],
  templateUrl: "./reset-password.html",
  styleUrl: "./reset-password.scss",
})
export class ResetPassword {
  private readonly accountApi = inject(AccountApi);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  readonly token = this.route.snapshot.queryParamMap.get("token") ?? "";
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);
  readonly isSubmitting = signal(false);

  readonly form = this.fb.group({
    password: ["", [Validators.required, Validators.minLength(8)]],
    confirmPassword: ["", Validators.required],
  });

  submit(): void {
    if (!this.token) {
      this.error.set("This password reset link is missing its token.");
      return;
    }

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    if (this.form.value.password !== this.form.value.confirmPassword) {
      this.error.set("Passwords do not match.");
      return;
    }

    this.error.set(null);
    this.success.set(null);
    this.isSubmitting.set(true);

    const req: ResetPasswordRequest = {
      token: this.token,
      password: this.form.value.password!,
    };

    this.accountApi.resetPassword(req).subscribe({
      next: () => {
        this.success.set("Password updated. Redirecting to login...");
        this.isSubmitting.set(false);
        void this.router.navigate([LOGIN_ROUTE]);
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
    if (error.status === 400) {
      this.error.set("This password reset link is invalid or expired.");
      return;
    }

    this.error.set("Failed to reset the password.");
  }
}
