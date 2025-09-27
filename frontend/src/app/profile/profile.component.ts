import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {
  user: any;
  errorMessage: string = '';

  constructor(private authService: AuthService) { }

  ngOnInit(): void {
    const user = this.authService.getUser();
    if (user && user.id) {
      this.authService.getProfile(user.id).subscribe({
        next: (profile) => {
          this.user = profile;
        },
        error: (error) => {
          this.errorMessage = error.error?.error || 'Failed to load profile.';
        }
      });
    } else {
      this.errorMessage = 'User not found.';
    }
  }
}
