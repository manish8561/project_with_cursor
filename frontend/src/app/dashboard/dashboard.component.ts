import { Component, OnInit } from '@angular/core';
import { RouterLink } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [RouterLink, CommonModule],
  template: `
    <div class="dashboard-container">
      <h1>Welcome to Dashboard</h1>
      <div *ngIf="user">
        <p>Hello, {{ user.name }}!</p>
        <p>Email: {{ user.email }}</p>
      </div>
      <div *ngIf="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      <a routerLink="/" class="btn btn-primary">Back to Home</a>
    </div>
  `,
  styles: [`
    .dashboard-container {
      max-width: 800px;
      margin: 2rem auto;
      padding: 2rem;
      text-align: center;
    }
    .error-message {
      color: red;
    }
  `]
})
export class DashboardComponent implements OnInit {
  user: any;
  errorMessage: string = '';

  constructor(private authService: AuthService) { }

  ngOnInit(): void {
    this.authService.getProfile().subscribe({
      next: (profile) => {
        this.user = profile;
      },
      error: (error) => {
        this.errorMessage = error.error?.error || 'Failed to load profile.';
      }
    });
  }
}
