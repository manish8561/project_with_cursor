import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  user: any;
  errorMessage: string = '';

  constructor(private authService: AuthService) { }

  ngOnInit(): void {
    if (!this.authService.getUserId()) {
      this.errorMessage = 'User not authenticated or user ID not available.';
      return;
    }

    this.authService.getProfile().subscribe({
      next: (profile) => {
        this.user = profile;
      },
      error: (error) => {
        console.error('Profile load error:', error);
        this.errorMessage = error.error?.error || 'Failed to load profile.';
      }
    });
  }
}
