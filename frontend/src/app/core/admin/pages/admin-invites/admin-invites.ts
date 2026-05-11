import { Component, inject, signal } from "@angular/core";
import { CommonModule } from "@angular/common";
import { HttpErrorResponse } from "@angular/common/http";
import { MatButtonModule } from "@angular/material/button";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { AdminAuth } from "../../services/admin-auth";
import { AdminInviteResponse } from "../../models/invite";

@Component({
  selector: "app-admin-invites",
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatFormFieldModule, MatInputModule],
  templateUrl: "./admin-invites.html",
  styleUrl: "./admin-invites.scss",
})
export class AdminInvites {
  private readonly adminAuth = inject(AdminAuth);

  readonly invite = signal<AdminInviteResponse | null>(null);
  readonly error = signal<string | null>(null);
  readonly isSubmitting = signal(false);
  readonly copyMessage = signal<string | null>(null);

  generateInvite(): void {
    this.error.set(null);
    this.copyMessage.set(null);
    this.isSubmitting.set(true);

    this.adminAuth.createInvite().subscribe({
      next: (invite) => {
        this.invite.set(invite);
        this.isSubmitting.set(false);
      },
      error: (error) => this.handleError(error),
    });
  }

  copyInviteLink(): void {
    const inviteUrl = this.invite()?.invite_url;
    if (!inviteUrl) return;

    void navigator.clipboard.writeText(inviteUrl).then(() => {
      this.copyMessage.set("Invite link copied.");
    });
  }

  private handleError(error: HttpErrorResponse): void {
    this.isSubmitting.set(false);
    if (error.status === 0) {
      this.error.set("Network error. Please check your connection.");
      return;
    }
    if (error.status === 401 || error.status === 403) {
      this.error.set("Your admin session is no longer valid.");
      return;
    }
    this.error.set("Failed to generate an invite link.");
  }
}
