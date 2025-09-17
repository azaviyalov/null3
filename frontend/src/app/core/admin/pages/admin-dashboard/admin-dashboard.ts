import { Component, inject, signal, OnInit } from "@angular/core";
import { Router } from "@angular/router";
import { MatTabsModule } from "@angular/material/tabs";
import { MatTableModule } from "@angular/material/table";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { CommonModule } from "@angular/common";
import { AdminAuth, User, RefreshToken } from "../../services/admin-auth";
import { UserDialogComponent } from "../user-dialog/user-dialog";

@Component({
  selector: "app-admin-dashboard",
  standalone: true,
  imports: [
    CommonModule,
    MatTabsModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule,
  ],
  templateUrl: "./admin-dashboard.html",
  styleUrl: "./admin-dashboard.scss",
})
export class AdminDashboard implements OnInit {
  private readonly adminAuth = inject(AdminAuth);
  private readonly router = inject(Router);
  private readonly dialog = inject(MatDialog);

  readonly users = signal<User[]>([]);
  readonly refreshTokens = signal<RefreshToken[]>([]);
  readonly isLoadingUsers = signal(false);
  readonly isLoadingTokens = signal(false);

  readonly displayedColumns: string[] = ["id", "name", "email", "actions"];
  readonly tokenDisplayedColumns: string[] = ["id", "user_id", "created_at", "expires_at", "actions"];

  ngOnInit() {
    this.loadUsers();
    this.loadRefreshTokens();
  }

  loadUsers() {
    this.isLoadingUsers.set(true);
    this.adminAuth.getUsers().subscribe({
      next: (users) => {
        this.users.set(users);
        this.isLoadingUsers.set(false);
      },
      error: (err) => {
        console.error("Failed to load users:", err);
        this.isLoadingUsers.set(false);
      },
    });
  }

  loadRefreshTokens() {
    this.isLoadingTokens.set(true);
    this.adminAuth.getRefreshTokens().subscribe({
      next: (tokens) => {
        this.refreshTokens.set(tokens);
        this.isLoadingTokens.set(false);
      },
      error: (err) => {
        console.error("Failed to load refresh tokens:", err);
        this.isLoadingTokens.set(false);
      },
    });
  }

  logout() {
    this.adminAuth.logout().subscribe({
      next: () => {
        this.router.navigate(["/admin/login"]);
      },
      error: (err) => {
        console.error("Failed to logout:", err);
        // Navigate anyway as logout might have worked on the server
        this.router.navigate(["/admin/login"]);
      },
    });
  }

  openCreateUserDialog() {
    const dialogRef = this.dialog.open(UserDialogComponent, {
      width: "500px",
      data: { mode: "create" },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.loadUsers();
      }
    });
  }

  editUser(user: User) {
    const dialogRef = this.dialog.open(UserDialogComponent, {
      width: "500px",
      data: { mode: "edit", user },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.loadUsers();
      }
    });
  }

  deleteUser(user: User) {
    if (confirm(`Are you sure you want to delete user "${user.name}"?`)) {
      this.adminAuth.deleteUser(user.id).subscribe({
        next: () => {
          this.loadUsers();
        },
        error: (err) => {
          console.error("Failed to delete user:", err);
          alert("Failed to delete user. Please try again.");
        },
      });
    }
  }

  deleteRefreshToken(token: RefreshToken) {
    if (confirm(`Are you sure you want to delete this refresh token?`)) {
      this.adminAuth.deleteRefreshToken(token.value).subscribe({
        next: () => {
          this.loadRefreshTokens();
        },
        error: (err) => {
          console.error("Failed to delete refresh token:", err);
          alert("Failed to delete refresh token. Please try again.");
        },
      });
    }
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString() + " " + date.toLocaleTimeString();
  }
}