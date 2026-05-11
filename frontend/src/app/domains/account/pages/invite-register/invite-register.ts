import { CommonModule } from "@angular/common";
import { Component, DestroyRef, inject, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { HttpErrorResponse } from "@angular/common/http";
import { ActivatedRoute, Router, RouterModule } from "@angular/router";
import { catchError, map, of, switchMap, tap } from "rxjs";
import { AccountApi } from "../../services/account-api";
import { InviteRegistrationRequest } from "../../models/invite";
import { UserSession } from "../../../session/services/user-session";

const ROOT_ROUTE = "/";
const LOGIN_PATTERN = /^[A-Za-z0-9_-]{3,32}$/;

@Component({
  selector: "app-invite-register",
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  templateUrl: "./invite-register.html",
  styleUrl: "./invite-register.scss",
})
export class InviteRegister {
  private readonly accountApi = inject(AccountApi);
  private readonly userSession = inject(UserSession);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly destroyRef = inject(DestroyRef);
  private readonly fb = inject(FormBuilder);

  readonly inviteToken = signal("");
  readonly inviteExpiresAt = signal<string | null>(null);
  readonly isLoadingInvite = signal(true);
  readonly error = signal<string | null>(null);
  readonly isSubmitting = signal(false);

  readonly form = this.fb.group({
    login: ["", [Validators.required, Validators.pattern(LOGIN_PATTERN)]],
    email: ["", [Validators.required, Validators.email]],
    password: ["", [Validators.required, Validators.minLength(8)]],
    confirmPassword: ["", Validators.required],
  });

  constructor() {
    this.route.paramMap
      .pipe(
        map((params) => params.get("token") ?? ""),
        tap((token) => {
          this.inviteToken.set(token);
          this.isLoadingInvite.set(true);
          this.error.set(null);
        }),
        switchMap((token) =>
          this.accountApi.getInvite(token).pipe(
            map((response) => response.expires_at),
            catchError((error: HttpErrorResponse) => {
              this.error.set(this.inviteLookupErrorMessage(error));
              return of(null);
            }),
          ),
        ),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe((expiresAt) => {
        this.isLoadingInvite.set(false);
        if (!expiresAt) {
          this.inviteExpiresAt.set(null);
          return;
        }
        this.inviteExpiresAt.set(expiresAt);
      });
  }

  submit(): void {
    if (!this.inviteToken()) {
      this.error.set("This invite link is invalid or expired.");
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
    this.isSubmitting.set(true);

    const req: InviteRegistrationRequest = {
      login: this.form.value.login!,
      email: this.form.value.email!,
      password: this.form.value.password!,
    };

    this.accountApi.registerWithInvite(this.inviteToken(), req).subscribe({
      next: (user) => {
        this.userSession.setAuthenticatedUser(user);
        this.isSubmitting.set(false);
        void this.router.navigate([ROOT_ROUTE]);
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
      this.error.set(
        this.httpErrorMessage(error, "Failed to complete registration."),
      );
      return;
    }
    if (error.status === 409) {
      this.error.set(
        this.httpErrorMessage(error, "That login or email is already in use."),
      );
      return;
    }

    this.error.set(
      this.httpErrorMessage(error, "Failed to complete registration."),
    );
  }

  private inviteLookupErrorMessage(error: HttpErrorResponse): string {
    if (error.status === 0) {
      return "Network error. Please check your connection.";
    }
    return this.httpErrorMessage(
      error,
      "This invite link is invalid or expired.",
    );
  }

  private httpErrorMessage(error: HttpErrorResponse, fallback: string): string {
    const message = error.error?.message;
    if (typeof message === "string" && message.trim().length > 0) {
      return message;
    }
    return fallback;
  }
}
